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

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("errror while reading env: %s", err.Error())
	}

	l, err := net.Listen("tcp", os.Getenv("SRV_PORT"))
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
