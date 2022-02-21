package game

type CreatureType int

const (
	CREATURETYPE_PLAYER CreatureType = iota
	CREATURETYPE_MONSTER
	CREATURETYPE_NPC
	CREATURETYPE_SUMMON_OWN
	CREATURETYPE_SUMMON_OTHERS
	CREATURETYPE_HIDDEN
)

type CreatureInterface interface {
	GetType() *CreatureType
}
