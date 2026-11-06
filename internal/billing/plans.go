package billing

type Plan struct { Code string `json:"code"`; MaxProjects int `json:"max_projects"`; MaxMembers int `json:"max_members"` }
var plans=map[string]Plan{
 "free":{Code:"free",MaxProjects:3,MaxMembers:2},
 "pro":{Code:"pro",MaxProjects:50,MaxMembers:10},
 "business":{Code:"business",MaxProjects:1000,MaxMembers:250},
}
func Get(code string)Plan{if p,ok:=plans[code];ok{return p};return plans["free"]}
