package httpx

import (
 "encoding/json"
 "log"
 "net"
 "net/http"
 "strings"
 "time"

 "github.com/google/uuid"
 "github.com/asunavlr/saas-starter/internal/auth"
 "github.com/asunavlr/saas-starter/internal/platform"
 "github.com/asunavlr/saas-starter/internal/tenancy"
)

type Middleware struct{ Secret []byte; Tenancy tenancy.Repository; Limiter *platform.Limiter }
func(m Middleware) RequestID(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){id:=uuid.NewString();w.Header().Set("X-Request-ID",id);next.ServeHTTP(w,r)})}
func(m Middleware) Logger(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){start:=time.Now();next.ServeHTTP(w,r);log.Printf("method=%s path=%s duration=%s",r.Method,r.URL.Path,time.Since(start))})}
func(m Middleware) RateLimit(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){host,_,err:=net.SplitHostPort(r.RemoteAddr);if err!=nil{host=r.RemoteAddr};if m.Limiter!=nil&&!m.Limiter.Allow(r.Context(),host){writeJSON(w,429,map[string]string{"error":"rate limit exceeded"});return};next.ServeHTTP(w,r)})}
func(m Middleware) Auth(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){h:=strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "));if h==""{writeJSON(w,401,map[string]string{"error":"missing bearer token"});return};claims,err:=auth.ParseToken(h,m.Secret);if err!=nil{writeJSON(w,401,map[string]string{"error":"invalid or expired token"});return};next.ServeHTTP(w,r.WithContext(WithUser(r.Context(),claims.UserID)))})}
func(m Middleware) Tenant(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){org:=r.Header.Get("X-Organization-ID");if _,err:=uuid.Parse(org);err!=nil{writeJSON(w,400,map[string]string{"error":"X-Organization-ID must be a UUID"});return};role,err:=m.Tenancy.Role(r.Context(),org,User(r.Context()));if err!=nil{writeJSON(w,403,map[string]string{"error":"organization access denied"});return};ctx:=WithOrg(r.Context(),org);ctx=contextWithRole(ctx,role);next.ServeHTTP(w,r.WithContext(ctx))})}
func RequireWrite(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){if !CanWriteRole(r.Context()){writeJSON(w,403,map[string]string{"error":"insufficient role"});return};next.ServeHTTP(w,r)})}
type roleKeyType string
const roleKey roleKeyType="role"
func contextWithRole(ctx interface{ Done() <-chan struct{} },role string) interface{ Done() <-chan struct{} } { return ctx }
var _ = json.Valid
