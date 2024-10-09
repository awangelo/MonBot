package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func ListGroups(client *whatsmeow.Client) {
	fmt.Println("Lista de Grupos:")
	groups, err := client.GetJoinedGroups()
	if err != nil {
		log.Fatalf("Erro ao buscar grupos: %v", err)
	}
	for _, group := range groups {
		fmt.Printf("Grupo: %s - JID: %s\n", group.Name, group.JID)
	}
}

func ListContacts(client *whatsmeow.Client) {
	contacts, err := client.Store.Contacts.GetAllContacts()
	if err != nil {
		log.Fatalf("Error fetching contacts: %v", err)
	}
	fmt.Println("Lista de Contatos e Grupos:")
	for jid, contact := range contacts {
		fmt.Printf("Contato: %s - JID: %s\n", contact.PushName, jid)
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

func IsBotMentioned(client *whatsmeow.Client, msg *events.Message) bool {
	// Verificar menções no ExtendedTextMessage
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		for _, jid := range msg.Message.ExtendedTextMessage.ContextInfo.MentionedJID {
			if jid == client.Store.ID.String() {
				return true
			}
		}
	}

	// Verificar menções no texto da mensagem
	text := msg.Message.GetConversation()
	if text == "" && msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	}
	return strings.Contains(text, "@"+client.Store.ID.User)
}

func PreventExit() {
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func InitLogger(msg *events.Message) {
	sender := msg.Info.Sender.String()

	// Verificar diferentes tipos de mensagens
	var text string
	if msg.Message.Conversation != nil {
		text = msg.Message.GetConversation()
	} else if msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	} else {
		// Outros tipos de mensagens (mídia, etc.)
		text = "media message"
	}

	log.Printf("Mensagem recebida no grupo de %s: %s\n", sender, text)
}
