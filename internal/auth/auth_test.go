package auth

import (
 "testing"
 "time"
)

func TestPasswordHashAndCheck(t *testing.T){
 hash,err:=HashPassword("a-strong-password");if err!=nil{t.Fatal(err)}
 if !CheckPassword(hash,"a-strong-password"){t.Fatal("expected password to match")}
 if CheckPassword(hash,"wrong-password"){t.Fatal("wrong password matched")}
}
func TestIssueAndParseToken(t *testing.T){
 secret:=[]byte("test-secret-with-at-least-32-characters")
 token,err:=IssueToken("user-123",secret,time.Minute);if err!=nil{t.Fatal(err)}
 claims,err:=ParseToken(token,secret);if err!=nil{t.Fatal(err)}
 if claims.UserID!="user-123"{t.Fatalf("got %q",claims.UserID)}
}
func TestRejectShortPassword(t *testing.T){if _,err:=HashPassword("short");err==nil{t.Fatal("expected short password error")}}
