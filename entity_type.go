package bonfire

import "math/rand"

type EntityType string

const (
	Character           = "character"
	Item_Consumable     = "item_consumable"
	Item_Key            = "item_key"
	Equipment_Weapon    = "equipment_weapon"
	Equipment_Armor     = "equipment_armor"
	Equipment_Accessory = "equipment_accessory"
	Location            = "location"
	Faction             = "faction"
	Event               = "event"
)

var AllEntityTypes = []EntityType{
	Character,
	Item_Consumable,
	Item_Key,
	Equipment_Weapon,
	Equipment_Armor,
	Equipment_Accessory,
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

func RandomEntityType() EntityType {
	entityIdx := rand.Intn(len(AllEntityTypes))
	return AllEntityTypes[entityIdx]
}
