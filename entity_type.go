package bonfire

import "math/rand"

type EntityType string

const (
	Character     = "character"
	Item          = "item"
	Relic         = "relic"
	Equipment     = "equipment"
	Location      = "location"
	Faction       = "faction"
	Event         = "event"
	World         = "world"
	EventCatalyst = "event_catalyst"
)

var AllEntityTypes = []EntityType{
	Character,
	Item,
	Relic,
	Equipment,
	Location,
	Faction,
	Event,
	World,
	EventCatalyst,
}

var AllNonSingletonEntityTypes = []EntityType{
	Character,
	Item,
	Relic,
	Equipment,
	Location,
	Faction,
	Event,
}

func IsValidEntityType(s string) bool {
	for _, et := range AllEntityTypes {
		if et == EntityType(s) {
			return true
		}
	}
	return false
}

func RandomNonSingletonEntityType() EntityType {
	entityIdx := rand.Intn(len(AllNonSingletonEntityTypes))
	return AllNonSingletonEntityTypes[entityIdx]
}
