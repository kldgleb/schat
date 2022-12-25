package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"schat"
	"time"

	"nhooyr.io/websocket"
)

const (
	subAddr = "ws://localhost:8080/subscribe"
	pubAddr = "http://localhost:8080/publish"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	msgCh := make(chan string, 16)
	subcribedCh := make(chan struct{})
	client := schat.NewWsClient(ctx, msgCh, "client")
	errc := make(chan error, 1)
	go func() {
		errc <- client.Subscribe(ctx, subAddr, &websocket.DialOptions{}, subcribedCh)
	}()
	<-subcribedCh

	client.Publish(ctx, pubAddr, "test")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		client.Shutdown()
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		client.Shutdown()
		log.Printf("terminating: %v", sig)
	}
}
