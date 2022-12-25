package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"schat"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	log.Printf("listening on http://%v", l.Addr())

	cs := schat.NewChatServer()
	s := &http.Server{
		Handler:      cs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}
