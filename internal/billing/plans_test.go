package billing

import "testing"

func TestPlanLimits(t *testing.T){
 if Get("free").MaxProjects!=3{t.Fatal("free plan should allow 3 projects")}
 if Get("pro").MaxProjects<=Get("free").MaxProjects{t.Fatal("pro should have higher limits")}
 if Get("unknown").Code!="free"{t.Fatal("unknown plans must fall back to free")}
}
