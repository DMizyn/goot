package network

import (
	"fmt"
	"github.com/rwxsu/goot/game"
	"github.com/rwxsu/goot/messages"
	"log"
	"net"
)

// Player message type
const (
	PlayerMessageTypeInfo   uint8 = 0x15
	PlayerMessageTypeCancel uint8 = 0x17
)

const (
	WALK_NORTH              uint8 = 0x65
	WALK_EAST               uint8 = 0x66
	WALK_SOUTH              uint8 = 0x67
	WALK_WEST               uint8 = 0x68
	PLAYER_SAY              uint8 = 0x96
	PLAYER_REQUEST_CHANNELS uint8 = 0x97
	OPEN_PRIVATE_CHANNEL    uint8 = 0x9A
	PARSE_FIGHT_MODES       uint8 = 0xa0
)

func ParseCommand(msg *messages.Message, player *game.Player, m *game.Map, code uint8) {
	switch code {
	case WALK_NORTH:
		if !SendMoveCreature(m, player, game.North, code) {
			SendSnapback(player)
		}
		return
	case WALK_EAST:
		if !SendMoveCreature(m, player, game.East, code) {
			SendSnapback(player)
		}
		return
	case WALK_SOUTH:
		if !SendMoveCreature(m, player, game.South, code) {
			SendSnapback(player)
		}
		return
	case WALK_WEST:
		if !SendMoveCreature(m, player, game.West, code) {
			SendSnapback(player)
		}
		return
	case PARSE_FIGHT_MODES:
		player.Tactic.FightMode = msg.ReadUint8()
		player.Tactic.ChaseOpponent = msg.ReadUint8()
		player.Tactic.AttackPlayers = msg.ReadUint8()
		return
	case PLAYER_SAY:
		if !parseSay(player, msg) {
			SendSnapback(player)
		}
	case PLAYER_REQUEST_CHANNELS:
		player.SendChannelsDialog()
		return
	case OPEN_PRIVATE_CHANNEL:
		receiver := msg.ReadString()
		log.Printf("Receiver: %s\n", receiver)
		playerOpenPrivateChannel(player, receiver)
		return
	default:
		log.Printf("Unknown code: %d\n", code)
		SendSnapback(player)
		return
	}
}

func playerOpenPrivateChannel(player *game.Player, receiver string) {

	if player == nil {
		return
	}

	/*	if game.GetPlayerByName(receiver) == nil {
		player.SendMessageWarning("A player with this name does not exist.")
		return
	}*/

	if player.Name == receiver {
		player.SendMessageWarning("You cannot set up a private message channel with yourself.")
		return
	}

	player.SendOpenPrivateChannel(receiver)

}

func parseSay(player *game.Player, msg *messages.Message) bool {
	var receiver string
	var channelId uint16
	var speakClass = game.SpeakClasses(msg.ReadUint8())

	switch speakClass {
	case game.TALKTYPE_PRIVATE_FROM:
		receiver = msg.ReadString()
		text := msg.ReadString()
		game.ChatManager.Messages <- game.Message{PlayerId: player.ID, ChannelId: channelId, SpeakClass: speakClass, Receiver: receiver, Text: text}
		break
	case game.TALKTYPE_PRIVATE_TO:
	case game.TALKTYPE_PRIVATE_RED_TO:
		receiver = msg.ReadString()
		channelId = 0
		break
	case game.TALKTYPE_CHANNEL_Y:
	case game.TALKTYPE_CHANNEL_R1:
		channelId = msg.ReadUint16()
		break
	default:
		channelId = 0
		break
	}

	var text = msg.ReadString()

	game.ChatManager.Messages <- game.Message{PlayerId: player.ID, ChannelId: channelId, SpeakClass: speakClass, Receiver: receiver, Text: text}
	return true
}

func SendInvalidClientVersion(c net.Conn) {
	msg := messages.NewMessage()
	msg.WriteUint8(0x0a)
	msg.WriteString("Only protocol 7.40 allowed!")
	SendMessage(c, msg)
}

func SendInvalidAccountOrPassword(c net.Conn) {
	msg := messages.NewMessage()
	msg.WriteUint8(0x0a)
	msg.WriteString("Please enter a valid account number and password.")
	SendMessage(c, msg)
}

func SendCharacterList(c net.Conn) {
	characters := make([]game.Creature, 2)
	characters[0].Name = "Goots"
	characters[0].World.Name = "test"
	characters[0].World.Port = 7171
	characters[1].Name = "GM Goots"
	characters[1].World.Name = "test"
	characters[1].World.Port = 7171
	res := messages.NewMessage()
	res.WriteUint8(0x14) // MOTD
	res.WriteString("Welcome to GoOT.")
	res.WriteUint8(0x64) // character list
	res.WriteUint8((uint8)(len(characters)))
	for i := 0; i < len(characters); i++ {
		res.WriteString(characters[i].Name)
		res.WriteString(characters[i].World.Name)
		res.WriteUint8(127)
		res.WriteUint8(0)
		res.WriteUint8(0)
		res.WriteUint8(1)
		res.WriteUint16(characters[i].World.Port)
	}
	res.WriteUint16(5) // premium days
	SendMessage(c, res)
}

func SendCharacterList2(c net.Conn, characters []game.Player) {
	res := messages.NewMessage()
	res.WriteUint8(0x14) // MOTD
	res.WriteString("Welcome to GoOT.")
	res.WriteUint8(0x64) // character list
	res.WriteUint8((uint8)(len(characters)))
	for i := 0; i < len(characters); i++ {
		res.WriteString(characters[i].Name)
		res.WriteString("World name")
		res.WriteUint8(127)
		res.WriteUint8(0)
		res.WriteUint8(0)
		res.WriteUint8(1)
		res.WriteUint16(7171)
	}
	res.WriteUint16(5) // premium days
	SendMessage(c, res)
}

func SendSnapback(player *game.Player) {
	msg := messages.NewMessage()
	msg.WriteUint8(0xb5)
	msg.WriteUint8(uint8(player.Direction))
	player.Client.WriteMessageToConn(msg)
	player.Client.SendMessageWarning("Sorry, not possible.")
}

func SendCancelMessage(c net.Conn, str string) {
	msg := messages.NewMessage()
	AddPlayerMessage(msg, str, PlayerMessageTypeCancel)
	SendMessage(c, msg)
}

func SendMoveCreature(m *game.Map, player *game.Player, direction game.DirectionType, code uint8) bool {
	var offset game.Offset
	var width, height uint16
	from := player.Position
	to := player.Position
	switch direction {
	case game.North:
		offset.X = -8
		offset.Y = -6
		width = 18
		height = 1
		to.Y--
		break
	case game.South:
		offset.X = -8
		offset.Y = 7
		width = 18
		height = 1
		to.Y++
		break
	case game.East:
		offset.X = 9
		offset.Y = -6
		width = 1
		height = 14
		to.X++
		break
	case game.West:
		offset.X = -8
		offset.Y = -6
		width = 1
		height = 14
		to.X--
		break
	}
	if !m.MoveCreature(&player.Creature, to, direction) {
		return false
	}
	msg := messages.NewMessage()
	msg.WriteUint8(0x6d)
	AddPosition(msg, from)
	msg.WriteUint8(0x01) // oldStackPos
	AddPosition(msg, to)
	msg.WriteUint8(code)
	msg.WriteUint16(0x63) // Creatureturn? In client's debug error.txt this is the "Parameter" field (0x63 == -1)
	AddMapDescription(msg, m, to, offset, width, height)
	player.Client.WriteMessageToConn(msg)
	return true
}

func SendAddCreature(character *game.Player, m *game.Map) {
	res := messages.NewMessage()
	res.WriteUint8(0x0a)
	res.WriteUint32(character.ID) // ID
	res.WriteUint16(0x32)         // ?
	// can report bugs?
	if character.Access > game.Regular {
		res.WriteUint8(0x01)
	} else {
		res.WriteUint8(0x00)
	}
	if character.Access >= game.Gamemaster {
		res.WriteUint8(0x0b)
		for i := 0; i < 32; i++ {
			res.WriteUint8(0xff)
		}
	}
	tile := m.GetTile(character.Position)
	tile.AddCreature(&character.Creature)
	res.WriteUint8(0x64)
	AddMapDescription(res, m, character.Position, game.Offset{X: -8, Y: -6, Z: 0}, 18, 14)
	AddMagicEffect(res, character.Position, 0x0a)
	AddInventory(res, character)
	AddStats(res, character)
	AddSkills(res, character)
	AddWorldLight(res, &character.World)
	AddCreatureLight(res, character)
	AddPlayerMessage(res, fmt.Sprintf("Welcome, %s.", character.Name), PlayerMessageTypeInfo)
	AddPlayerMessage(res, "TODO: Last Login String 01-01-1970", PlayerMessageTypeInfo)
	AddCreatureLight(res, character)
	AddIcons(res, character)
	character.Client.WriteMessageToConn(res)
}

func AddCreatureLight(msg *messages.Message, c *game.Player) {
	msg.WriteUint8(0x8d)
	msg.WriteUint32(c.ID)
	msg.WriteUint8(c.Light.Level)
	msg.WriteUint8(c.Light.Color)
}

func AddWorldLight(msg *messages.Message, w *game.World) {
	msg.WriteUint8(0x82)
	msg.WriteUint8(w.Light.Level)
	msg.WriteUint8(w.Light.Color)
}

func AddIcons(msg *messages.Message, c *game.Player) {
	msg.WriteUint8(0xa2)
	msg.WriteUint8(c.Icons)
}

func AddSkills(msg *messages.Message, c *game.Player) {
	msg.WriteUint8(0xa1)
	msg.WriteUint8(c.Fist.Level)
	msg.WriteUint8(c.Fist.Percent)
	msg.WriteUint8(c.Club.Level)
	msg.WriteUint8(c.Club.Percent)
	msg.WriteUint8(c.Sword.Level)
	msg.WriteUint8(c.Sword.Percent)
	msg.WriteUint8(c.Axe.Level)
	msg.WriteUint8(c.Axe.Percent)
	msg.WriteUint8(c.Distance.Level)
	msg.WriteUint8(c.Distance.Percent)
	msg.WriteUint8(c.Shielding.Level)
	msg.WriteUint8(c.Shielding.Percent)
	msg.WriteUint8(c.Fishing.Level)
	msg.WriteUint8(c.Fishing.Percent)
}

func AddStats(msg *messages.Message, c *game.Player) {
	msg.WriteUint8(0xa0) // send player stats
	msg.WriteUint16(c.HealthNow)
	msg.WriteUint16(c.HealthMax)
	msg.WriteUint16(c.Cap)
	msg.WriteUint32(c.Combat.Experience)
	msg.WriteUint8(c.Combat.Level)
	msg.WriteUint8(c.Combat.Percent)
	msg.WriteUint16(c.ManaNow)
	msg.WriteUint16(c.ManaMax)
	msg.WriteUint8(c.Magic.Level)
	msg.WriteUint8(c.Magic.Percent)
}

func AddInventory(msg *messages.Message, c *game.Player) {
	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotHead)

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotNecklace)

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotBackpack)
	msg.WriteUint16(0x7c4) // backpack

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotBody)
	msg.WriteUint16(0x9a8) // magic plate armor

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotShield)

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotWeapon)
	msg.WriteUint16(0x997) // crossbow

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotLegs)

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotFeet)

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotRing)

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotAmmo)
	msg.WriteUint16(0x9ef) // bolts
	msg.WriteUint8(33)     // count
}

func AddMapDescription(msg *messages.Message, m *game.Map, pos game.Position, offset game.Offset, width, height uint16) {
	AddPosition(msg, pos)
	pos.Offset(offset)
	skip := uint8(0)
	if pos.Z < 8 {
		for z := (int8)(7); z > -1; z-- {
			for x := (uint16)(0); x < width; x++ {
				for y := (uint16)(0); y < height; y++ {
					tile := m.GetTile(game.Position{X: pos.X + x, Y: pos.Y + y, Z: (uint8)(z)})
					if tile != nil {
						if skip > 0 {
							msg.WriteUint8(skip - 1)
							msg.WriteUint8(0xff)
						}
						skip = 1
						AddTile(msg, tile)
					} else if skip == 0xfe {
						msg.WriteUint8(skip)
						msg.WriteUint8(0xff)
						skip = 0
					} else {
						skip++
					}
				}
			}
		}
	} else { // TODO: underground

	}
	// Remainder
	if skip > 0 {
		msg.WriteUint8(skip - 1)
		msg.WriteUint8(0xff)
	}
}

func AddPosition(msg *messages.Message, pos game.Position) {
	msg.WriteUint16(pos.X)
	msg.WriteUint16(pos.Y)
	msg.WriteUint8(pos.Z)
}

func AddMagicEffect(msg *messages.Message, pos game.Position, kind uint8) {
	msg.WriteUint8(0x83)
	AddPosition(msg, pos)
	msg.WriteUint8(kind)
}

func AddCreature(msg *messages.Message, c *game.Creature) {
	msg.WriteUint16(0x61) // unknown creature
	msg.WriteUint32(0x00) // something about caching known creatures
	msg.WriteUint32(c.ID)
	msg.WriteString(c.Name)
	msg.WriteUint8((uint8)(c.HealthNow*100/c.HealthMax) + 1)
	msg.WriteUint8(uint8(c.Direction)) // look dir
	msg.WriteUint8(c.Outfit.Type)
	msg.WriteUint8(c.Outfit.Head)
	msg.WriteUint8(c.Outfit.Body)
	msg.WriteUint8(c.Outfit.Legs)
	msg.WriteUint8(c.Outfit.Feet)
	msg.WriteUint8(c.Light.Level)
	msg.WriteUint8(c.Light.Color)
	msg.WriteUint16(c.Speed)
	msg.WriteUint8(c.Skull)
	msg.WriteUint8(c.Party)
}

// AddTile adds all the tile items and creatures WITHOUT the end of tile
//delimeter (0xSKIP-0xff)
func AddTile(msg *messages.Message, tile *game.Tile) {
	for _, i := range tile.Items {
		msg.WriteUint16(i.ID)
	}
	for _, c := range tile.Creatures {
		AddCreature(msg, c)
	}
}

func AddPlayerMessage(msg *messages.Message, str string, kind uint8) {
	msg.WriteUint8(0xb4)
	msg.WriteUint8(kind)
	msg.WriteString(str)
}
