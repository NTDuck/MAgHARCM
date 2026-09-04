package agents

import "testing"

func TestBlackboardScheduler(t *testing.T) {
	bb := NewBlackboard()
	bb.Register("u1")
	wu, ok := bb.Claim("u1", "analyzer")
	if !ok || wu == nil {
		t.Fatal("claim failed")
	}
	if wu.Status != "CLAIMED" {
		t.Fatalf("status: got %s want CLAIMED", wu.Status)
	}
	if !bb.Complete("u1", "ok") {
		t.Fatal("complete failed")
	}
	evs, units := bb.Snapshot()
	if len(evs) < 3 {
		t.Fatalf("events: got %d want >=3", len(evs))
	}
	if len(bb.AuditLog) == 0 {
		t.Fatal("audit log empty")
	}
	u, ok := units["u1"]
	if !ok || u.Status != "COMPLETED" {
		t.Fatalf("unit status: got %v want COMPLETED", u)
	}
}
