package httpx

import "context"

type contextKey string
const userKey contextKey="user_id"
const orgKey contextKey="organization_id"
func WithUser(ctx context.Context,id string)context.Context{return context.WithValue(ctx,userKey,id)}
func User(ctx context.Context)string{v,_:=ctx.Value(userKey).(string);return v}
func WithOrg(ctx context.Context,id string)context.Context{return context.WithValue(ctx,orgKey,id)}
func Org(ctx context.Context)string{v,_:=ctx.Value(orgKey).(string);return v}
