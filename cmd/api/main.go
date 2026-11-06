package main

import (
 "context"
 "log"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"

 "github.com/go-chi/chi/v5"
 "github.com/asunavlr/saas-starter/internal/audit"
 "github.com/asunavlr/saas-starter/internal/config"
 httpx "github.com/asunavlr/saas-starter/internal/httpx"
 "github.com/asunavlr/saas-starter/internal/platform"
 "github.com/asunavlr/saas-starter/internal/projects"
 "github.com/asunavlr/saas-starter/internal/tenancy"
)

func main(){
 cfg,err:=config.Load();if err!=nil{log.Fatal(err)}
 ctx,cancel:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM);defer cancel()
 db,err:=platform.OpenPostgres(ctx,cfg.DatabaseURL);if err!=nil{log.Fatalf("database: %v",err)};defer db.Close()
 limiter:=platform.NewLimiter(cfg.RedisAddr,cfg.RedisPassword,cfg.RateLimitRPM);defer limiter.Close()
 mw:=httpx.Middleware{Secret:[]byte(cfg.JWTSecret),Tenancy:tenancy.Repository{DB:db},Limiter:limiter}
 h:=httpx.Handlers{DB:db,Secret:[]byte(cfg.JWTSecret),TokenTTL:cfg.AccessTokenTTL,Projects:projects.Repository{DB:db},Audit:audit.Repository{DB:db}}
 r:=chi.NewRouter();r.Use(mw.RequestID,mw.Logger,mw.RateLimit)
 r.Get("/healthz",func(w http.ResponseWriter,r *http.Request){w.WriteHeader(http.StatusOK);w.Write([]byte("ok"))})
 r.Get("/readyz",func(w http.ResponseWriter,r *http.Request){if err:=db.PingContext(r.Context());err!=nil{http.Error(w,"database unavailable",503);return};w.WriteHeader(200);w.Write([]byte("ready"))})
 r.Route("/v1",func(v chi.Router){
  v.Post("/auth/register",h.Register);v.Post("/auth/login",h.Login)
  v.Group(func(a chi.Router){a.Use(mw.Auth);a.Get("/me",h.Me)
   a.Group(func(t chi.Router){t.Use(mw.Tenant);t.Get("/projects",h.ListProjects);t.With(httpx.RequireWrite).Post("/projects",h.CreateProject);t.With(httpx.RequireWrite).Delete("/projects/{id}",h.DeleteProject)})
  })
 })
 srv:=&http.Server{Addr:cfg.HTTPAddr,Handler:r,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:30*time.Second,IdleTimeout:60*time.Second}
 go func(){log.Printf("API listening on %s",cfg.HTTPAddr);if err:=srv.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{log.Fatalf("http server: %v",err)}}()
 <-ctx.Done();shutdown,cancelShutdown:=context.WithTimeout(context.Background(),10*time.Second);defer cancelShutdown();if err:=srv.Shutdown(shutdown);err!=nil{log.Printf("shutdown: %v",err)}
}
