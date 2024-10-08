package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
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
	client := whatsmeow.NewClient(device, nil)
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

	listGroups(client)

	groupJID := "556182100810-1580939047@g.us"
	for i := range 100 {
		err := sendMessageToGroup(client, groupJID, fmt.Sprintf("Mensagem %d", i))
		if err != nil {
			log.Fatalf("Erro ao enviar mensagem: %v", err)
		}
		time.Sleep(time.Second * 2)
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func sendMessageToGroup(client *whatsmeow.Client, groupJID string, message string) error {
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

func listGroups(client *whatsmeow.Client) {
	fmt.Println("Lista de Grupos:")
	groups, err := client.GetJoinedGroups()
	if err != nil {
		log.Fatalf("Erro ao buscar grupos: %v", err)
	}

	for _, group := range groups {
		fmt.Printf("Grupo: %s - JID: %s\n", group.Name, group.JID)
	}
}
