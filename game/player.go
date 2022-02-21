package game

const debug = true

type Player struct {
	Creature
	Client
}

func (p Player) GetPosition() *Position {
	return &p.Position
}

func (p Player) GetPlayer() *Player {
	return &p
}
