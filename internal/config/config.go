package config

import (
 "fmt"
 "os"
 "strconv"
 "time"
)

type Config struct {
 Env, HTTPAddr, DatabaseURL, RedisAddr, RedisPassword, JWTSecret string
 AccessTokenTTL time.Duration
 RateLimitRPM int
}

func Load() (Config, error) {
 c := Config{
  Env: get("APP_ENV","development"),
  HTTPAddr: get("HTTP_ADDR",":8080"),
  DatabaseURL: os.Getenv("DATABASE_URL"),
  RedisAddr: get("REDIS_ADDR","localhost:6379"),
  RedisPassword: os.Getenv("REDIS_PASSWORD"),
  JWTSecret: os.Getenv("JWT_SECRET"),
  RateLimitRPM: 120,
 }
 if c.DatabaseURL == "" { return c, fmt.Errorf("DATABASE_URL is required") }
 if len(c.JWTSecret) < 32 { return c, fmt.Errorf("JWT_SECRET must contain at least 32 characters") }
 ttl, err := time.ParseDuration(get("ACCESS_TOKEN_TTL","15m")); if err != nil { return c, fmt.Errorf("ACCESS_TOKEN_TTL: %w",err) }; c.AccessTokenTTL=ttl
 rpm, err := strconv.Atoi(get("RATE_LIMIT_RPM","120")); if err != nil || rpm < 1 { return c, fmt.Errorf("RATE_LIMIT_RPM must be a positive integer") }; c.RateLimitRPM=rpm
 return c,nil
}
func get(k, fallback string) string { if v:=os.Getenv(k); v!="" { return v }; return fallback }
