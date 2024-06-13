package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	nc, err := nats.Connect("0.0.0.0:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error getting JetStream context: %v", err)
	}

	// Check if the stream already exists
	streamName := "orders"
	streamInfo, err := js.StreamInfo(streamName)
	if err != nil {
		// Stream does not exist, create a new one
		if errors.Is(err, nats.ErrStreamNotFound) {
			_, err = js.AddStream(&nats.StreamConfig{
				Name:     streamName,
				Subjects: []string{"orders"},
			})
			if err != nil {
				log.Fatalf("Error adding stream: %v", err)
			}
			log.Println("Stream created successfully.")
		} else {
			log.Fatalf("Error checking stream info: %v", err)
		}
	} else {
		// Stream exists, log the info
		log.Printf("Stream %s already exists. Created at: %s\n", streamName, streamInfo.Created.Format(time.DateTime))
	}

	log.Println(Green + "Ready to publish to JetStream." + Reset)

	for {
		fmt.Println("Generating a new random order")
		fmt.Print("Enter order_uid: ")
		// Wait for input
		if scanner.Scan() {
			input := scanner.Text()
			randOrder, err := json.Marshal(GenerateRandomOrder(input))
			_, err = js.Publish("orders", randOrder)
			log.Printf("A new order with %suid=%s%s published to JetStream.", Green, input, Reset)
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}
		} else {
			// Handle errors or EOF
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
			break
		}
	}
}
