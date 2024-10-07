package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

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
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
