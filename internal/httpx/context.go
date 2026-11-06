package httpx

import "context"

type contextKey string
const userKey contextKey="user_id"
const orgKey contextKey="organization_id"
const roleKey contextKey="role"
func WithUser(ctx context.Context,id string)context.Context{return context.WithValue(ctx,userKey,id)}
func User(ctx context.Context)string{v,_:=ctx.Value(userKey).(string);return v}
func WithOrg(ctx context.Context,id string)context.Context{return context.WithValue(ctx,orgKey,id)}
func Org(ctx context.Context)string{v,_:=ctx.Value(orgKey).(string);return v}
func WithRole(ctx context.Context,role string)context.Context{return context.WithValue(ctx,roleKey,role)}
func Role(ctx context.Context)string{v,_:=ctx.Value(roleKey).(string);return v}
func CanWriteRole(ctx context.Context)bool{return Role(ctx)=="owner"||Role(ctx)=="admin"}
