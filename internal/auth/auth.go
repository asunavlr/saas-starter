package auth

import (
 "errors"
 "time"

 "github.com/golang-jwt/jwt/v5"
 "golang.org/x/crypto/bcrypt"
 "github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Claims struct {
 UserID string `json:"uid"`
 jwt.RegisteredClaims
}

func HashPassword(password string) (string,error) {
 if len(password)<10 { return "", errors.New("password must be at least 10 characters") }
 b,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost); return string(b),err
}
func CheckPassword(hash,password string) bool { return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))==nil }
func IssueToken(userID string, secret []byte, ttl time.Duration) (string,error) {
 claims:=Claims{UserID:userID,RegisteredClaims:jwt.RegisteredClaims{Subject:userID,ID:uuid.NewString(),IssuedAt:jwt.NewNumericDate(time.Now()),ExpiresAt:jwt.NewNumericDate(time.Now().Add(ttl))}}
 return jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString(secret)
}
func ParseToken(token string,secret []byte)(*Claims,error) {
 parsed,err:=jwt.ParseWithClaims(token,&Claims{},func(t *jwt.Token)(any,error){ if t.Method!=jwt.SigningMethodHS256{return nil,errors.New("unexpected signing method")}; return secret,nil })
 if err!=nil{return nil,err}; claims,ok:=parsed.Claims.(*Claims); if !ok||!parsed.Valid{return nil,errors.New("invalid token")}; return claims,nil
}
