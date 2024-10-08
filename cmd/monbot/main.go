package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var client *whatsmeow.Client
var groupJID types.JID

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	// Jid do grupo deve ser guardado em .env
	groupJIDString := os.Getenv("GROUP_JID")
	if groupJIDString == "" {
		log.Fatal("GROUP_JID is not defined in .env")
	}

	// Parse do JID
	var parseErr error
	groupJID, parseErr = types.ParseJID(groupJIDString)
	if parseErr != nil {
		log.Fatalf("Error parsing GROUP_JID: %v", parseErr)
	}
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		handleGroupMessage(v)
	}
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

	if isBotMentioned(msg) {
		fmt.Println("O bot foi mencionado!")
		replyToMention(msg)
	}
}

func main() {
	// Container envolve um storage sqlite
	container, err := sqlstore.New("sqlite3", "file:whatsmeow.db?_foreign_keys=true", nil)
	if err != nil {
		log.Fatalf("Error creating container: %v", err)
	}

	// Criar um novo dispositivo unico e client
	device, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Error getting first device: %v", err)
	}
	client = whatsmeow.NewClient(device, nil)
	client.AddEventHandler(eventHandler)

	// Solicitar um codigo QR caso nao tenha
	if client.Store.ID == nil {
		// qrChan eh um channel que gera um novo QRCode sempre que o anterior expira
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
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
		err = client.Connect()
		if err != nil {
			log.Fatalf("Error connecting client: %v", err)
		}
		log.Println("connected")
	}

	listGroups()
	// listContacts()

	for i := range 2 {
		err := sendMessageToGroup(groupJID.String(), fmt.Sprintf("Mensagem %d", i))
		if err != nil {
			log.Fatalf("Erro ao enviar mensagem: %v", err)
		}
		time.Sleep(time.Second * 1)
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func isBotMentioned(msg *events.Message) bool {
	if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.ContextInfo != nil {
		mentionedJIDs := msg.Message.ExtendedTextMessage.ContextInfo.MentionedJID
		for _, jid := range mentionedJIDs {
			if jid == client.Store.ID.String() {
				return true
			}
		}
	}
	// Verificar também no texto da mensagem (para casos de menção sem formatação)
	text := msg.Message.GetConversation()
	if text == "" && msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	}
	return strings.Contains(text, "@"+client.Store.ID.User)
}

func replyToMention(msg *events.Message) {
	response := fmt.Sprintf("Olá @%s, você me mencionou!", msg.Info.PushName)

	replyMsg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(response),
			ContextInfo: &waProto.ContextInfo{
				StanzaID:      &msg.Info.ID,
				Participant:   proto.String(msg.Info.Sender.String()),
				QuotedMessage: msg.Message,
				MentionedJID:  []string{msg.Info.Sender.String()},
			},
		},
	}

	_, err := client.SendMessage(context.Background(), msg.Info.Chat, replyMsg)
	if err != nil {
		log.Printf("Erro ao responder à menção: %v", err)
	} else {
		fmt.Println("Resposta enviada com sucesso!")
	}
}

func sendMessageToGroup(groupJID string, message string) error {
	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("JID inválido: %v", err)
	}

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err = client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return fmt.Errorf("Erro ao enviar mensagem: %v", err)
	}

	fmt.Println("Mensagem enviada com sucesso!")
	return nil
}

func listGroups() {
	fmt.Println("Lista de Grupos:")
	groups, err := client.GetJoinedGroups()
	if err != nil {
		log.Fatalf("Erro ao buscar grupos: %v", err)
	}

	for _, group := range groups {
		fmt.Printf("Grupo: %s - JID: %s\n", group.Name, group.JID)
	}
}

func listContacts() {
	// Obter todos os contatos conhecidos pelo dispositivo
	contacts, err := client.Store.Contacts.GetAllContacts()
	if err != nil {
		log.Fatalf("Error fetching contacts: %v", err)
	}
	fmt.Println("Lista de Contatos e Grupos:")
	for jid, contact := range contacts {
		fmt.Printf("Contato: %s - JID: %s\n", contact.PushName, jid)
	}
	fmt.Println("done")
}
