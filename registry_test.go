package registry

import "testing"

func TestNew(t *testing.T) {
	r := New()

	if r == nil {
		t.Fatal("expected registry")
	}

	if r.jobs == nil {
		t.Error("jobs map not initialized")
	}

	if r.scheduler == nil {
		t.Error("scheduler not initialized")
	}
}
