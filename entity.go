package bonfire

import (
	"fmt"
	"regexp"
	"time"
)

type Entity struct {
	Name      string     `json:"name"`
	Id        string     `json:"id"`
	Type      EntityType `json:"type"`
	ShortDesc string     `json:"short_desc"`
	LongDesc  string     `json:"long_desc"`
	CreatedAt time.Time  `json:"created_at"`
}

type EntityReferenceHint struct {
	Id        string     `json:"id"`
	Type      EntityType `json:"type"`
	ShortDesc string     `json:"short_desc"`
}

type UnknownReference struct {
	Id                  string
	ReferencingEntityId string
}

func (e *Entity) validate() error {
	if !IsValidEntityType(string(e.Type)) {
		return fmt.Errorf("invalid entity type: %s", e.Type)
	}
	return nil
}

func (e *Entity) ParseReferences(replacer func(fullMatch, id, innerText string) string) string {
	re := regexp.MustCompile(`<ref id=["']([^"']+)["']>([^<]+)</ref>`)
	return re.ReplaceAllStringFunc(e.LongDesc, func(fullMatch string) string {
		matches := re.FindStringSubmatch(fullMatch)
		if len(matches) != 3 {
			return fullMatch
		}
		id := matches[1]
		innerText := matches[2]
		return replacer(fullMatch, id, innerText)
	})
}
