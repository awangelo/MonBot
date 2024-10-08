package handler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/awangelo/MonBot/internal/bot"
	"github.com/awangelo/MonBot/internal/utils"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

var client *whatsmeow.Client

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		handleGroupMessage(v)
	}
}

func InitHandler(c *whatsmeow.Client) {
	c.AddEventHandler(eventHandler)
	client = c
}

func handleGroupMessage(msg *events.Message) {
	// Ignorar se não for mensagem de um grupo
	if msg.Info.Chat.Server != "g.us" {
		return
	}

	sender := msg.Info.Sender.String()
	var text string

	// Verificar diferentes tipos de mensagens
	if msg.Message.Conversation != nil {
		text = msg.Message.GetConversation()
	} else if msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	} else {
		// Outros tipos de mensagens (mídia, etc.)
		text = "Mensagem não textual"
	}

	fmt.Printf("Mensagem recebida no grupo de %s: %s\n", sender, text)

	if utils.IsBotMentioned(client, msg) {
		fmt.Println("O bot foi mencionado!")
		bot.ReplyToMention(client, msg)
	}
}

func LoginByQr(client *whatsmeow.Client) {
	// Solicitar um codigo QR caso nao tenha
	if client.Store.ID == nil {
		// qrChan eh um channel que gera um novo QRCode sempre que o anterior expira
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			log.Fatalf("Error connecting client: %v", err)
		}
		// Caso receber um QRCode, printa no term usando qrterminal
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Se a sessão já estiver logada, conecte-se diretamente
		err := client.Connect()
		if err != nil {
			log.Fatalf("Error connecting client: %v", err)
		}
		log.Println("connected")
	}
}
