package game

import "time"

const debug = true

type Player struct {
	Creature
	Client
	AccountNumber uint32
	AccountType   uint16
	PremiumEndsAt time.Time
	Sex           uint8
	Vocation      uint8
	LastLogin     time.Time
	LastLogout    time.Time
	LastIp        string
	SkullTime     time.Time
	TownId        uint8
}

func (p Player) GetPosition() *Position {
	return &p.Position
}

func (p Player) GetPlayer() *Player {
	return &p
}
