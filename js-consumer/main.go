package main

import (
	"bytes"
	"fmt"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"net/http"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"

func main() {
	// Connect to NATS server on localhost
	nc, err := nats.Connect("0.0.0.0:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()
	log.Println(Green + "Подключение к NATS выполнено успешно" + Reset)

	url := "http://localhost:8080/api/orders"
	client := &http.Client{}

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error getting JetStream context: %v", err)
	}

	// Subscribe to the stream
	subj := "orders"
	sub, err := js.Subscribe(subj, func(m *nats.Msg) {
		log.Printf("%sReceived%s data of subject \"%s%s%s\". Sending a POST request to %s ", Green, Reset, Green, subj, Reset, url)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(m.Data))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			log.Printf("%sError sending request%s: %v\n", Red, Reset, err)
		} else {
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("%sError reading response body%s: %v", Red, Reset, err)
			}
			fmt.Printf("Response status: %s\n", resp.Status)
			fmt.Printf("Response body: %s\n", body)
		}

		//log.Printf("Received a message: %s\n", string(m.Data))
		m.Ack()
	}, nats.Durable("my-durable"), nats.ManualAck())
	if err != nil {
		log.Fatalf("Error subscribing to JetStream: %v", err)
	}
	defer sub.Unsubscribe()

	// Keep the connection alive
	select {}
}
