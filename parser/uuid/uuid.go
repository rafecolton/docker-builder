package uuid

import gouuid "github.com/nu7hatch/gouuid"

/*
UUID is currently a wrapper for the gouuid library with some extra functions
*/
type UUID gouuid.UUID

type seededUUIDGenerator struct {
	UUID
}

/*
A Generator generates uuids either randomly or the same one every time (for
test purposes)
*/
type Generator interface {
	NextUUID() (string, error)
}

/*
NewUUIDGenerator returns a UUIDGenerator.  If passed (true), the generator will
produce a unique uuid every time.  If passed (false), the generator will
produce the same uuid every time.
*/
func NewUUIDGenerator(random bool) Generator {
	if random {
		return &randomUUIDGenerator{}
	}

	return &seededUUIDGenerator{}
}

type randomUUIDGenerator struct {
	UUID
}

func (gen *seededUUIDGenerator) NextUUID() (string, error) {
	u, err := gouuid.NewV5(gouuid.NamespaceURL, []byte("0"))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (gen *randomUUIDGenerator) NextUUID() (string, error) {
	u, err := gouuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
