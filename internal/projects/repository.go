package projects

import (
 "context"
 "database/sql"
 "errors"

 "github.com/google/uuid"
)

type Project struct { ID string `json:"id"`; OrganizationID string `json:"organization_id"`; Name string `json:"name"`; Description string `json:"description"`; CreatedAt string `json:"created_at"` }
type Repository struct{ DB *sql.DB }
func(r Repository) List(ctx context.Context,org string)([]Project,error){
 rows,err:=r.DB.QueryContext(ctx,`SELECT id,organization_id,name,description,created_at::text FROM projects WHERE organization_id=$1 ORDER BY created_at DESC LIMIT 100`,org);if err!=nil{return nil,err};defer rows.Close()
 out:=[]Project{};for rows.Next(){var p Project;if err:=rows.Scan(&p.ID,&p.OrganizationID,&p.Name,&p.Description,&p.CreatedAt);err!=nil{return nil,err};out=append(out,p)};return out,rows.Err()
}
func(r Repository) Create(ctx context.Context,org,user,name,description string)(Project,error){
 var p Project
 err:=r.DB.QueryRowContext(ctx,`INSERT INTO projects(organization_id,created_by,name,description) VALUES($1,$2,$3,$4) RETURNING id,organization_id,name,description,created_at::text`,org,user,name,description).Scan(&p.ID,&p.OrganizationID,&p.Name,&p.Description,&p.CreatedAt)
 return p,err
}
func(r Repository) Delete(ctx context.Context,org,id string)error{
 if _,err:=uuid.Parse(id);err!=nil{return errors.New("invalid project id")}
 res,err:=r.DB.ExecContext(ctx,`DELETE FROM projects WHERE organization_id=$1 AND id=$2`,org,id);if err!=nil{return err};n,_:=res.RowsAffected();if n==0{return sql.ErrNoRows};return nil
}
