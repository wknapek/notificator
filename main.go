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
	messageFile, err := os.Open(*messagePath)
	defer func(messageFile *os.File) {
		err = messageFile.Close()
		if err != nil {

		}
	}(messageFile)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(messageFile)
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second * time.Duration(*intervals))
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("starting sending messages on ", *url)
	messages := make(chan string, 10)
	stopSending := make(chan bool, 1)
	go readLines(scanner, messages, stopSending, signalChan)
	for {
		select {
		case <-ticker.C:
			if len(messages) == 0 {
				return
			}
			notificator.SendMessages(ctx, url, messages)
		case <-stopSending:
			fmt.Println("\nReceived an interrupt, stopping graceful shutdown...")
			cancel()
			ticker.Stop()
			close(stopSending)
			close(messages)
			return
		}
	}
}

func readLines(scanner *bufio.Scanner, messages chan string, stop chan bool, signalChan chan os.Signal) {
	for scanner.Scan() {
		select {
		case <-signalChan:
			fmt.Println("\nReceived an interrupt, stopping graceful shutdown...")
			stop <- true
			return
		default:
			messages <- scanner.Text()
		}
	}
}
