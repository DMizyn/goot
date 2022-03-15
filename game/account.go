package game

type Account struct {
	Id          uint32
	Password    string
	AccountType uint8
	Characters  []Player
}
