package game

import (
	"log"
)

const MINIMUM_LEVEL_TO_SEND_PRIVATE = 1

type Manager struct {
	normalChannels map[int32]*chatChannel
	Messages       chan Message
}

func NewChat() *Manager {
	return &Manager{
		normalChannels: make(map[int32]*chatChannel),
		Messages:       make(chan Message, 10),
	}
}

func (manager Manager) ProceedData() {
	for message := range manager.Messages {

		player := GetPlayerByID(message.PlayerId)

		if player == nil {
			return
		}

		switch message.SpeakClass {
		case TALKTYPE_PRIVATE_FROM:
			playerSpeakTo(player, message.SpeakClass, message.Receiver, message.Text)
			break
		case TALKTYPE_SAY:
			log.Printf("PlayerId: %d, ChannelId: %d, SpeakClass: %d, Receiver: %s, Text: %s", message.PlayerId, message.ChannelId, message.SpeakClass, message.Receiver, message.Text)
			player.Client.SendMessageWarning(message.Text)
			break
		default:
			log.Printf("PlayerId: %d, ChannelId: %d, SpeakClass: %d, Receiver: %s, Text: %s", message.PlayerId, message.ChannelId, message.SpeakClass, message.Receiver, message.Text)
		}
	}
}

func playerSpeakTo(player *Player, class SpeakClasses, receiverName string, text string) {
	receiver := GetPlayerByName(receiverName)

	if receiver == nil {
		player.SendMessageWarning("A player with this name is not online.")
	}

	/*	if (type == TALKTYPE_PRIVATE_RED_TO && (player->hasFlag(PlayerFlag_CanTalkRedPrivate) || player->getAccountType() >= ACCOUNT_TYPE_GAMEMASTER)) {
			type = TALKTYPE_PRIVATE_RED_FROM;
		} else {
			type = TALKTYPE_PRIVATE_FROM;
		}*/

	if player.Combat.Level < MINIMUM_LEVEL_TO_SEND_PRIVATE {
		player.Client.SendMessageWarning("You may not send private messages unless you have reached level " + string(rune(MINIMUM_LEVEL_TO_SEND_PRIVATE)) + ".")
	}
}

type Message struct {
	PlayerId   uint32
	ChannelId  uint16
	SpeakClass SpeakClasses
	Receiver   string
	Text       string
}
