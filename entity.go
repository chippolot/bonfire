package bonfire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Entity struct {
	Name      string
	Id        string
	Type      EntityType
	Lore      string
	CreatedAt time.Time
}

func (e *Entity) validate() error {
	if !IsValidEntityType(string(e.Type)) {
		return fmt.Errorf("invalid entity type: %s", e.Type)
	}
	return nil
}

func (t *Entity) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
