package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notify/notificator"
)

func main() {
	url := flag.String("url", "", "url address to send notification")
	intervals := flag.Uint("intervals", 0, "intervals seconds to send notification")
	messagePath := flag.String("messages", "", "file with message")
	flag.Parse()
	var messages []string
	messageFile, err := os.Open(*messagePath)
	defer messageFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(messageFile)
	for scanner.Scan() {
		msg := scanner.Text()
		messages = append(messages, msg)
	}
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second * time.Duration(*intervals))
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("starting sending messages on ", *url)
	for {
		select {
		case <-ticker.C:
			notificator.SendMessages(ctx, *url, messages)
		case <-signalChan:
			fmt.Println("\nReceived an interrupt, stopping graceful shutdown...")
			cancel()
			ticker.Stop()
			return
		}
	}
}
