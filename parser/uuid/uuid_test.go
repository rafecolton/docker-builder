package uuid_test

import (
	. "github.com/rafecolton/docker-builder/parser/uuid"
	"testing"
)

func TestRandomGenerator(t *testing.T) {
	var random = NewUUIDGenerator(true)
	alpha, _ := random.NextUUID()
	beta, _ := random.NextUUID()
	if alpha == beta {
		t.Errorf("expected random uuids to be different")
	}
}

func TestSeededGenerator(t *testing.T) {
	var seeded = NewUUIDGenerator(false)
	alpha, _ := seeded.NextUUID()
	beta, _ := seeded.NextUUID()
	if alpha != beta {
		t.Errorf("expected uuids from seeded generator to be identical")
	}
}
