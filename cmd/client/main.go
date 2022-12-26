package main

import (
	"log"
	"schat"
	"schat/cli"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	subAddr = "ws://localhost:8080/subscribe"
	pubAddr = "http://localhost:8080/publish"
)

func main() {
	readyToChat := make(chan struct{})
	msgCh := make(chan schat.MsgForm, 16)
	sendCh := make(chan schat.MsgForm)
	m := cli.NewMainModel(readyToChat, msgCh, sendCh)
	p := tea.NewProgram(m)
	go func() {
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	<-readyToChat
	for {
		msg := <-sendCh
		var newMsg cli.NewMsg
		p.Send(newMsg)
		msgCh <- msg
	}
}

// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	defer cancel()
// 	msgCh := make(chan string, 16)
// 	subcribedCh := make(chan struct{})
// 	client := schat.NewWsClient(ctx, msgCh, "client")
// 	errc := make(chan error, 1)
// 	go func() {
// 		errc <- client.Subscribe(ctx, subAddr, &websocket.DialOptions{}, subcribedCh)
// 	}()
// 	<-subcribedCh

// 	client.Publish(ctx, pubAddr, "test")

// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, os.Interrupt)
// 	select {
// 	case err := <-errc:
// 		client.Shutdown()
// 		log.Printf("failed to serve: %v", err)
// 	case sig := <-sigs:
// 		client.Shutdown()
// 		log.Printf("terminating: %v", sig)
// 	}
