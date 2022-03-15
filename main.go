package main

import (
	"github.com/rwxsu/goot/database"
	_ "github.com/rwxsu/goot/database"
	"github.com/rwxsu/goot/game"
	"github.com/rwxsu/goot/messages"
	"github.com/rwxsu/goot/network"
	"log"
	"net"
	"path/filepath"
)

func main() {
	m := make(game.Map)

	const sectors = "data/map/sectors/*"
	filenames, _ := filepath.Glob(sectors)

	for _, fn := range filenames {
		m.LoadSector(fn)
	}

	l, err := net.Listen("tcp", ":7171")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	acceptConnections(l, m)
}

func acceptConnections(l net.Listener, m game.Map) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handleConnection(c, m)
	}
}

func handleConnection(c net.Conn, m game.Map) {
	client := game.Client{
		Client: c,
	}
	creature := game.Creature{}
	player := &game.Player{Creature: creature, Client: client}
connectionLoop:
	for {
		req := network.RecvMessage(c)
		if req == nil {
			return
		}
		code := req.ReadUint8()
		switch code {
		case 0x01: // request character list
			parseFirstPacket(c, req)
			break connectionLoop
		case 0x0a: // request character login
			req.SkipBytes(2) // os := req.ReadUint16()
			if req.ReadUint16() != 740 {
				network.SendInvalidClientVersion(c)
				break connectionLoop
			}
			req.SkipBytes(5)
			characterName := req.ReadString()
			log.Printf("PlayerName: %s", characterName)
			err := database.LoadPlayerByName(player, characterName)
			if err != nil {
				log.Println(err)
			}
			game.AddPlayer(player)
			network.SendAddCreature(player, &m)
		case 0x14: //
			break connectionLoop
		default:
			network.ParseCommand(req, player, &m, code)
		}
	}
	if err := c.Close(); err != nil {
		log.Printf("Unable to close connection %v\n", err)
	}
}

func parseFirstPacket(c net.Conn, msg *messages.Message) bool {
	msg.SkipBytes(2) // os := req.ReadUint16()
	if msg.ReadUint16() != 740 {
		network.SendInvalidClientVersion(c)
		return false
	}

	msg.SkipBytes(12)

	accNumber := msg.ReadUint32()
	password := msg.ReadString()

	log.Printf("SCID_6476465276 Account number: %d, password: %s ", accNumber, password)
	account, error := database.GetAccountById(accNumber)
	if error != nil {
		log.Fatalln(error)
		network.SendInvalidAccountOrPassword(c)
		return false
	}

	if account.Id != accNumber || account.Password != password {
		network.SendInvalidAccountOrPassword(c)
		return false
	}

	network.SendCharacterList(c, account.Characters)
	return true
}
