package handler

import (
	"log"

	"github.com/awangelo/MonBot/internal/bot"
	"github.com/awangelo/MonBot/internal/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	client *whatsmeow.Client
	group  types.JID
)

func handleGroupMessage(msg *events.Message) {
	// Ignorar se nao for mensagem do grupo especifico
	if msg.Info.Chat != group {
		return
	}

	if bot.IsBotMentioned(client, msg) {
		log.Println("O bot foi mencionado!")
		bot.ReplyToMention(client, msg)
		return
	}

	if msg.Message.GetConversation()[0] == '!' {
		HandleCommand(client, msg)
		return
	}
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		utils.LogMessageEvent(v)
		handleGroupMessage(v)
	}
}

// InitHandler inicializa o handler de eventos
func InitHandler(c *whatsmeow.Client, g types.JID) {
	c.AddEventHandler(eventHandler)
	client = c
	group = g
}
