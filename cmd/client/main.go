package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"schat"
	"schat/cli"

	tea "github.com/charmbracelet/bubbletea"
	"nhooyr.io/websocket"
)

const (
	subAddr = "ws://localhost:8080/subscribe"
	pubAddr = "http://localhost:8080/publish"
)

func main() {
	readyToChat := make(chan struct{})
	chatMsgCh := make(chan schat.MsgForm)
	sendCh := make(chan schat.MsgForm)
	m := cli.NewMainModel(readyToChat, chatMsgCh, sendCh)
	p := tea.NewProgram(m)
	go func() {
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	<-readyToChat

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subcribedCh := make(chan struct{})
	msgCh := make(chan schat.MsgForm, 16)
	client := schat.NewWsClient(ctx, msgCh, "client")
	errc := make(chan error, 1)

	go func() {
		errc <- client.Subscribe(ctx, subAddr, &websocket.DialOptions{}, subcribedCh)
	}()
	<-subcribedCh

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
App:
	for {
		select {
		case err := <-errc:
			client.Shutdown()
			log.Printf("failed to serve: %v", err)
			break App
		case sig := <-sigs:
			client.Shutdown()
			log.Printf("terminating: %v", sig)
			break App
		case msg := <-sendCh:
			client.Publish(ctx, pubAddr, msg)
		case msg := <-msgCh:
			var newMsg cli.NewMsg
			p.Send(newMsg)
			chatMsgCh <- msg
		}
	}
}
