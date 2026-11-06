package httpx

import (
 "database/sql"
 "encoding/json"
 "errors"
 "net/http"
 "strings"
 "time"

 "github.com/asunavlr/saas-starter/internal/audit"
 "github.com/asunavlr/saas-starter/internal/auth"
 "github.com/asunavlr/saas-starter/internal/billing"
 "github.com/asunavlr/saas-starter/internal/projects"
 "github.com/google/uuid"
 "github.com/go-chi/chi/v5"
)

type Handlers struct { DB *sql.DB; Secret []byte; TokenTTL time.Duration; Projects projects.Repository; Audit audit.Repository }
type credentials struct { Name string `json:"name"`; Email string `json:"email"`; Password string `json:"password"` }
func(h Handlers) Register(w http.ResponseWriter,r *http.Request){
 var in credentials;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid JSON"});return}
 in.Name=strings.TrimSpace(in.Name);in.Email=strings.ToLower(strings.TrimSpace(in.Email))
 if in.Name==""||len(in.Name)>120||!strings.Contains(in.Email,"@"){writeJSON(w,400,map[string]string{"error":"valid name and email are required"});return}
 hash,err:=auth.HashPassword(in.Password);if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return}
 tx,err:=h.DB.BeginTx(r.Context(),nil);if err!=nil{writeJSON(w,500,map[string]string{"error":"could not start transaction"});return};defer tx.Rollback()
 var userID,orgID string
 err=tx.QueryRowContext(r.Context(),`INSERT INTO users(name,email,password_hash) VALUES($1,$2,$3) RETURNING id`,in.Name,in.Email,hash).Scan(&userID)
 if err!=nil{writeJSON(w,409,map[string]string{"error":"email already registered"});return}
 slug:=strings.ReplaceAll(strings.ToLower(in.Name)," ","-")+"-"+uuid.NewString()[:8]
 err=tx.QueryRowContext(r.Context(),`INSERT INTO organizations(name,slug) VALUES($1,$2) RETURNING id`,in.Name+"'s workspace",slug).Scan(&orgID)
 if err==nil{_,err=tx.ExecContext(r.Context(),`INSERT INTO memberships(organization_id,user_id,role) VALUES($1,$2,'owner')`,orgID,userID)}
 if err==nil{_,err=tx.ExecContext(r.Context(),`INSERT INTO subscriptions(organization_id,plan,status) VALUES($1,'free','active')`,orgID)}
 if err!=nil{writeJSON(w,500,map[string]string{"error":"could not create workspace"});return}
 if err=tx.Commit();err!=nil{writeJSON(w,500,map[string]string{"error":"could not commit registration"});return}
 token,err:=auth.IssueToken(userID,h.Secret,h.TokenTTL);if err!=nil{writeJSON(w,500,map[string]string{"error":"could not issue token"});return}
 writeJSON(w,201,map[string]any{"access_token":token,"token_type":"Bearer","user":map[string]string{"id":userID,"name":in.Name,"email":in.Email},"organization":map[string]string{"id":orgID,"name":in.Name+"'s workspace","role":"owner"},"plan":"free"})
}
func(h Handlers) Login(w http.ResponseWriter,r *http.Request){
 var in credentials;if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid JSON"});return}
 var id,name,email,hash string
 err:=h.DB.QueryRowContext(r.Context(),`SELECT id,name,email,password_hash FROM users WHERE email=$1`,strings.ToLower(strings.TrimSpace(in.Email))).Scan(&id,&name,&email,&hash)
 if err!=nil||!auth.CheckPassword(hash,in.Password){writeJSON(w,401,map[string]string{"error":"invalid email or password"});return}
 token,err:=auth.IssueToken(id,h.Secret,h.TokenTTL);if err!=nil{writeJSON(w,500,map[string]string{"error":"could not issue token"});return}
 rows,err:=h.DB.QueryContext(r.Context(),`SELECT o.id,o.name,m.role,o.plan FROM organizations o JOIN memberships m ON m.organization_id=o.id WHERE m.user_id=$1 ORDER BY o.created_at`,id)
 if err!=nil{writeJSON(w,500,map[string]string{"error":"could not load organizations"});return};defer rows.Close()
 orgs:=[]map[string]string{};for rows.Next(){var oid,oname,role,plan string;if rows.Scan(&oid,&oname,&role,&plan)==nil{orgs=append(orgs,map[string]string{"id":oid,"name":oname,"role":role,"plan":plan})}}
 writeJSON(w,200,map[string]any{"access_token":token,"token_type":"Bearer","user":map[string]string{"id":id,"name":name,"email":email},"organizations":orgs})
}
func(h Handlers) ListProjects(w http.ResponseWriter,r *http.Request){items,err:=h.Projects.List(r.Context(),Org(r.Context()));if err!=nil{writeJSON(w,500,map[string]string{"error":"could not list projects"});return};writeJSON(w,200,map[string]any{"data":items})}
func(h Handlers) CreateProject(w http.ResponseWriter,r *http.Request){
 var in struct{Name string `json:"name"`;Description string `json:"description"`};if err:=json.NewDecoder(r.Body).Decode(&in);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid JSON"});return}
 in.Name=strings.TrimSpace(in.Name);if in.Name==""||len(in.Name)>160{writeJSON(w,400,map[string]string{"error":"project name must contain 1-160 characters"});return}
 var plan string;if err:=h.DB.QueryRowContext(r.Context(),`SELECT plan FROM organizations WHERE id=$1`,Org(r.Context())).Scan(&plan);err!=nil{writeJSON(w,500,map[string]string{"error":"could not load plan"});return}
 var count int;if err:=h.DB.QueryRowContext(r.Context(),`SELECT count(*) FROM projects WHERE organization_id=$1`,Org(r.Context())).Scan(&count);err!=nil{writeJSON(w,500,map[string]string{"error":"could not check project quota"});return}
 if count>=billing.Get(plan).MaxProjects{writeJSON(w,402,map[string]string{"error":"project limit reached for current plan"});return}
 p,err:=h.Projects.Create(r.Context(),Org(r.Context()),User(r.Context()),in.Name,in.Description);if err!=nil{writeJSON(w,500,map[string]string{"error":"could not create project"});return}
 _=h.Audit.Record(r.Context(),Org(r.Context()),User(r.Context()),"project.created","project",p.ID,map[string]any{"name":p.Name})
 writeJSON(w,201,map[string]any{"data":p})
}
func(h Handlers) DeleteProject(w http.ResponseWriter,r *http.Request){id:=chi.URLParam(r,"id");err:=h.Projects.Delete(r.Context(),Org(r.Context()),id);if errors.Is(err,sql.ErrNoRows){writeJSON(w,404,map[string]string{"error":"project not found"});return};if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};_=h.Audit.Record(r.Context(),Org(r.Context()),User(r.Context()),"project.deleted","project",id,nil);w.WriteHeader(http.StatusNoContent)}
func(h Handlers) Me(w http.ResponseWriter,r *http.Request){var name,email string;if err:=h.DB.QueryRowContext(r.Context(),`SELECT name,email FROM users WHERE id=$1`,User(r.Context())).Scan(&name,&email);err!=nil{writeJSON(w,404,map[string]string{"error":"user not found"});return};writeJSON(w,200,map[string]string{"id":User(r.Context()),"name":name,"email":email})}
func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_=json.NewEncoder(w).Encode(v)}
