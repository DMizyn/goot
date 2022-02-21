package game

import (
	"github.com/rwxsu/goot/messages"
	"net"
)

const (
	MessageNone = iota
	MessageSay
	MessageWhisper
	MessageYell
	MessagePrivateFrom
	MessagePrivateTo
	MessageChannelManagement
	MessageChannel
	MessageChannelHighlight
	MessageSpell
	MessageNpcFrom
	MessageNpcTo
	MessageGamemasterBroadcast
	MessageGamemasterChannel
	MessageGamemasterPrivateFrom
	MessageGamemasterPrivateTo
	MessageLogin
	MessageWarning
	MessageGame
	MessageFailure
	MessageLook
	MessageDamageDealed
	MessageDamageReceived
	MessageHeal
	MessageExp
	MessageDamageOthers
	MessageHealOthers
	MessageExpOthers
	MessageStatus
	MessageLoot
	MessageTradeNpc
	MessageGuild
	MessagePartyManagement
	MessageParty
	MessageBarkLow
	MessageBarkLoud
	MessageReport
	MessageHotkeyUse
	MessageTutorialHint
	MessageThankyou
	MessageMarket
	MessageMana
	MessageBeyondLast

	// deprecated
	MessageMonsterYell
	MessageMonsterSay
	MessageRed
	MessageBlue
	MessageRVRChannel
	MessageRVRAnswer
	MessageRVRContinue
	MessageGameHighlight
	MessageNpcFromStartBlock
	LastMessage
	MessageInvalid = 255
)

type Client struct {
	Client net.Conn
}

func (c Client) SendChannelsDialog() {
	msg := messages.NewMessage()
	msg.WriteUint8(0xAB)
	msg.WriteUint8(1)
	msg.WriteUint16(123)
	msg.WriteString("DUPA")
	c.WriteMessageToConn(msg)
}

func (c Client) SendOpenPrivateChannel(receiver string) {
	msg := messages.NewMessage()
	msg.WriteUint8(0xAD)
	msg.WriteString(receiver)
	c.WriteMessageToConn(msg)
}

func (c Client) SendMessageSay(text string) {
	c.sendMessage(text, MessageSay)
}

func (c Client) SendMessagePrivateFrom(text string) {
	c.sendMessage(text, MessagePrivateFrom)
}

func (c Client) SendMessagePrivateTo(text string) {
	c.sendMessage(text, MessagePrivateTo)
}

func (c Client) SendMessageWarning(text string) {
	c.sendMessage(text, MessageWarning)
}

func (c Client) sendMessage(text string, kind uint8) {
	msg := messages.NewMessage()
	addPlayerMessage(msg, text, kind)
	c.WriteMessageToConn(msg)
}

func addPlayerMessage(msg *messages.Message, str string, kind uint8) {
	msg.WriteUint8(0xb4)
	msg.WriteUint8(kind)
	msg.WriteString(str)
}

func (c Client) WriteMessageToConn(msg *messages.Message) {
	c.Client.Write(msg.Buffer)
	if debug {
		msg.HexDump("send")
	}
}
