package audit

import (
 "context"
 "database/sql"
 "encoding/json"
)

type Repository struct{ DB *sql.DB }
func(r Repository) Record(ctx context.Context,orgID,userID,action,resourceType,resourceID string,metadata map[string]any)error{
 raw,err:=json.Marshal(metadata);if err!=nil{return err}
 _,err=r.DB.ExecContext(ctx,`INSERT INTO audit_events(organization_id,actor_user_id,action,resource_type,resource_id,metadata) VALUES($1,$2,$3,$4,$5,$6)`,orgID,userID,action,resourceType,resourceID,raw)
 return err
}
