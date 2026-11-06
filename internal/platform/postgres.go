package platform

import (
 "context"
 "database/sql"
 "time"

 _ "github.com/jackc/pgx/v5/stdlib"
)

func OpenPostgres(ctx context.Context, dsn string) (*sql.DB,error) {
 db,err:=sql.Open("pgx",dsn); if err!=nil{return nil,err}
 db.SetMaxOpenConns(25); db.SetMaxIdleConns(5); db.SetConnMaxLifetime(30*time.Minute)
 pingCtx,cancel:=context.WithTimeout(ctx,5*time.Second); defer cancel()
 if err=db.PingContext(pingCtx); err!=nil { db.Close(); return nil,err }; return db,nil
}
