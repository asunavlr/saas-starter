package platform

import (
 "context"
 "sync"
 "time"

 "github.com/redis/go-redis/v9"
)

type Limiter struct { client *redis.Client; mu sync.Mutex; local map[string]bucket; limit int }
type bucket struct { count int; reset time.Time }
func NewLimiter(addr,password string,limit int)*Limiter {
 return &Limiter{client:redis.NewClient(&redis.Options{Addr:addr,Password:password}),local:map[string]bucket{},limit:limit}
}
func(l *Limiter) Allow(ctx context.Context,key string) bool {
 if l.client!=nil {
  n,err:=l.client.Incr(ctx,"rl:"+key).Result()
  if err==nil { if n==1 { _=l.client.Expire(ctx,"rl:"+key,time.Minute).Err() }; return n<=int64(l.limit) }
 }
 l.mu.Lock(); defer l.mu.Unlock(); now:=time.Now(); b:=l.local[key]
 if now.After(b.reset) { b=bucket{reset:now.Add(time.Minute)} }
 b.count++; l.local[key]=b; return b.count<=l.limit
}
func(l *Limiter) Close()error { if l.client==nil{return nil}; return l.client.Close() }
