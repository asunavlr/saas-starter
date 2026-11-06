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

func(r Repository) List(ctx context.Context,orgID string)([]map[string]any,error){
 rows,err:=r.DB.QueryContext(ctx,`SELECT id,actor_user_id,action,resource_type,resource_id,metadata,created_at::text FROM audit_events WHERE organization_id=$1 ORDER BY created_at DESC LIMIT 100`,orgID)
 if err!=nil{return nil,err};defer rows.Close()
 items:=[]map[string]any{}
 for rows.Next(){
  var id int64;var actor sql.NullString;var action,resourceType,resourceID,created string;var raw []byte
  if err:=rows.Scan(&id,&actor,&action,&resourceType,&resourceID,&raw,&created);err!=nil{return nil,err}
  var metadata map[string]any;_ = json.Unmarshal(raw,&metadata)
  items=append(items,map[string]any{"id":id,"actor_user_id":actor.String,"action":action,"resource_type":resourceType,"resource_id":resourceID,"metadata":metadata,"created_at":created})
 }
 return items,rows.Err()
}
