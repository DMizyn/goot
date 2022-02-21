package game

// Access
const (
	Regular uint8 = iota + 1
	Tutor
	SeniorTutor
	Gamemaster
	God
)

// Slot
const (
	SlotHead uint8 = iota + 1
	SlotNecklace
	SlotBackpack
	SlotBody
	SlotShield // Right hand
	SlotWeapon // Left hand
	SlotLegs
	SlotRing
	SlotFeet
	SlotAmmo
)

// Slot Capacity
const (
	SlotEmpty    uint8 = 0x79
	SlotNotEmpty uint8 = 0x78
)

// Tactic contains the user selected options right of equipment slots in client
type Tactic struct {
	// Fight modes
	// 1 - Offensive fighting
	// 2 - Balanced fighting
	// 3 - Defensive fighting
	FightMode uint8

	// 0 - Stand while fighting
	// 1 - Chase opponent
	ChaseOpponent uint8

	// 0 - You cannot attack unmarked players
	// 1 - You can attack any player
	AttackPlayers uint8
}

type Skill struct {
	Experience uint32 // tries
	Level      uint8
	Percent    uint8
}

type Outfit struct {
	Type uint8
	Head uint8
	Body uint8
	Legs uint8
	Feet uint8
}

type Creature struct {
	Position     Position
	Access       uint8
	World        World
	ID           uint32
	creatureType CreatureType
	Name         string
	Direction    DirectionType
	Outfit       Outfit
	Speed        uint16
	Skull        uint8
	Party        uint8
	Cap          uint16
	HealthNow    uint16
	HealthMax    uint16
	ManaNow      uint16
	ManaMax      uint16
	Combat       Skill
	Magic        Skill
	Fist         Skill
	Club         Skill
	Sword        Skill
	Axe          Skill
	Distance     Skill
	Shielding    Skill
	Fishing      Skill
	Light        Light
	Icons        uint8
	Tactic       Tactic
}

func (c Creature) GetType() *CreatureType {
	return &c.creatureType
}
