package tenancy

import (
 "context"
 "database/sql"
)

type Repository struct { DB *sql.DB }
func(r Repository) Role(ctx context.Context,orgID,userID string)(string,error) {
 var role string
 err:=r.DB.QueryRowContext(ctx,`SELECT role FROM memberships WHERE organization_id=$1 AND user_id=$2`,orgID,userID).Scan(&role)
 return role,err
}
func CanWrite(role string)bool{return role=="owner"||role=="admin"}
