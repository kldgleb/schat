package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

const (
	TextType = iota + 1
	BinaryType
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/subscribe", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "Internal error")

	err = publish()
	if err != nil {
		panic(err)
	}

	for {
		_, str, err := c.Read(ctx)
		if err != nil {
			panic(err)
		}
		log.Printf("received: %s", str)
	}

	c.Close(websocket.StatusNormalClosure, "")
}

func publish() error {
	// values := map[string]string{"msg": "test"}
	// json_data, err := json.Marshal(values)
	// if err != nil {
	// 	return err
	// }
	resp, err := http.Post("http://localhost:8080/publish", "application/json", bytes.NewBuffer([]byte("test")))
	if err != nil {
		return err
	}
	fmt.Println("Publish status code: ", resp.StatusCode)
	return nil
}
