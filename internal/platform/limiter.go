package platform

import (
 "context"
 "sync"
 "time"

 "github.com/redis/go-redis/v9"
)

// Limiter uses Redis for a shared fixed-window counter. If Redis is unavailable,
// it falls back to a process-local counter suitable for development only.
type Limiter struct {
 client *redis.Client
 rpm int
 mu sync.Mutex
 buckets map[string]bucket
}
type bucket struct { count int; expires time.Time }

func NewLimiter(addr,password string,rpm int)*Limiter {
 if rpm<1 { rpm=120 }
 client:=redis.NewClient(&redis.Options{Addr:addr,Password:password})
 ctx,cancel:=context.WithTimeout(context.Background(),750*time.Millisecond)
 defer cancel()
 if err:=client.Ping(ctx).Err();err!=nil { _=client.Close();client=nil }
 return &Limiter{client:client,rpm:rpm,buckets:make(map[string]bucket)}
}
func(l *Limiter) Allow(ctx context.Context,key string)bool {
 if l==nil{return true}
 window:=time.Minute
 if l.client!=nil {
  n,err:=l.client.Incr(ctx,"saas:rate:"+key).Result()
  if err==nil {
   if n==1 { _=l.client.Expire(ctx,"saas:rate:"+key,window).Err() }
   return n<=int64(l.rpm)
  }
 }
 now:=time.Now();l.mu.Lock();defer l.mu.Unlock()
 b:=l.buckets[key];if now.After(b.expires){b=bucket{expires:now.Add(window)}}
 b.count++;l.buckets[key]=b
 if len(l.buckets)>10000 { for k,v:=range l.buckets { if now.After(v.expires){delete(l.buckets,k)} } }
 return b.count<=l.rpm
}
func(l *Limiter) Close()error { if l==nil||l.client==nil{return nil};return l.client.Close() }
