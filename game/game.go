package game

type DirectionType int

var ChatManager *Manager
var players map[uint32]*Player

func init() {
	ChatManager = NewChat()
	players = make(map[uint32]*Player)

	go ChatManager.ProceedData()
}

func GetPlayerByName(name string) *Player {
	for _, v := range players {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func GetPlayerByID(id uint32) *Player {
	return players[id]
}

func AddPlayer(player *Player) {
	players[player.ID] = player
}

func RemovePlayer(player *Player) {
	delete(players, player.ID)
}

// Direction
const (
	North DirectionType = iota
	East
	South
	West
)

// Position is the real in-game position
type Position struct {
	X uint16
	Y uint16
	Z uint8
}

type Offset struct {
	X, Y, Z int8
}

func (pos *Position) Offset(offset Offset) {
	pos.X += (uint16)(offset.X)
	pos.Y += (uint16)(offset.Y)
	pos.Z += (uint8)(offset.Z)
}

// Light has the same structure for both creatures and world
type Light struct {
	Level uint8
	Color uint8
}
