package handler

import (
	"log"

	"github.com/awangelo/MonBot/internal/bot"
	"github.com/awangelo/MonBot/internal/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var client *whatsmeow.Client

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go utils.InitLogger(v)
		handleGroupMessage(v)
	}
}

func InitHandler(c *whatsmeow.Client) {
	c.AddEventHandler(eventHandler)
	client = c
}

func handleGroupMessage(msg *events.Message) {
	// Ignorar se não for mensagem de um grupo e se não for uma mensagem de texto
	if msg.Info.Chat.Server != "g.us" {
		return
	}

	if utils.IsBotMentioned(client, msg) {
		log.Println("O bot foi mencionado!")
		bot.ReplyToMention(client, msg)
		return
	}

	if msg.Message.GetConversation()[0] == '!' {
		bot.HandleCommand(client, msg)
		return
	}
}
