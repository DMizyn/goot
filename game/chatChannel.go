package game

type chatChannel struct {
	int           int32
	name          string
	users         map[int32]*Creature
	publicChannel bool
}
