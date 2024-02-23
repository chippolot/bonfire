package bonfire

import "fmt"

type Entity struct {
	Name string
	Id   string
	Type EntityType
	Lore string
}

func (e *Entity) validate() error {
	if !IsValidEntityType(string(e.Type)) {
		return fmt.Errorf("invalid entity type: %s", e.Type)
	}
	return nil
}
