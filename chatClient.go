package schat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

type wsClient struct {
	ctx   context.Context
	msgCh chan MsgForm
	name  string
}

type MsgForm struct {
	Msg  string `json:"msg"`
	Name string `json:"name"`
}


func NewWsClient(ctx context.Context, msgCh chan MsgForm, name string) *wsClient {
	return &wsClient{
		ctx:   ctx,
		msgCh: msgCh,
		name:  name,
	}
}

func (c *wsClient) Subscribe(ctx context.Context, addr string, opts *websocket.DialOptions, subcribedCh chan<- struct{}) error {
	conn, _, err := websocket.Dial(ctx, addr, opts)
	if err != nil {
		close(subcribedCh)
		return err
	}
	subcribedCh <- struct{}{}
	close(subcribedCh)
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	for {
		var msgForm MsgForm
		_, data, err := conn.Read(ctx)
		if err != nil {
			return err
		}
		decrypted := decrypt(key, data)
		err = json.Unmarshal([]byte(decrypted), &msgForm)
		if err != nil {
			return err
		}

		select {
		case c.msgCh <- msgForm:
			continue
		case <-ctx.Done():
			return nil
		}

	}
}

func (c *wsClient) Publish(ctx context.Context, addr string, msg MsgForm) error {
	json_data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	json_data = encrypt(key, json_data)

	client := http.Client{}
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("wrong status code: 202 != %d /n", resp.StatusCode)
	}
	return nil
}

func (c *wsClient) Shutdown() {
	close(c.msgCh)
	fmt.Println("client was shutted down")
}
