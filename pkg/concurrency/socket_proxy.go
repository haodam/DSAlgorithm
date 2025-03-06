package main

//
//import (
//	"context"
//	"fmt"
//	"strings"
//	"time"
//
//	fiberws "github.com/gofiber/websocket/v2"
//	"golang.org/x/net/websocket"
//	"golang.org/x/sync/errgroup"
//)
//
//type (
//	messageTransportIf interface {
//		ReadMessage() (int, []byte, error)
//		WriteMessage(int, []byte) error
//		Close() error
//	}
//
//	connectMessage struct {
//		AuthToken string `json:"auth_token"`
//	}
//
//	webSocketNatsWrapper struct {
//		*websocket.Conn
//	}
//
//	webSocketClient struct {
//		IsConnected bool
//		*fiberws.Conn
//	}
//
//	webSocketProxy struct {
//		client *webSocketClient
//		server *webSocketNatsWrapper
//	}
//)
//
//func (w *webSocketClient) ReadMessage() (int, []byte, error) {
//	if !w.IsConnected {
//		// blocking
//		msgType, msg, err := w.Conn.ReadMessage()
//		if err != nil {
//			return msgType, msg, err
//		}
//		stringArray := strings.Split(string(msg), "CONNECT ")
//		if len(stringArray) == 2 {
//			connectMessage, err := common.TransformToType[connectMessage]([]byte(stringArray[1]))
//			if err != nil {
//				return msgType, msg, fmt.Errorf("invalid connect message: %s", err)
//			}
//
//			_, errApp := helpers.VerifySignedToken(
//				config.TokenSecret,
//				config.TokenEncryptSecret,
//				connectMessage.AuthToken,
//			)
//			if errApp != nil {
//				return msgType, msg, fmt.Errorf("invalid verify token: %s", errApp)
//			}
//
//			w.IsConnected = true
//		}
//
//		return msgType, msg, err
//	}
//	return w.Conn.ReadMessage()
//}
//
//func (w *webSocketNatsWrapper) ReadMessage() (int, []byte, error) {
//	var msg []byte
//
//	err := websocket.Message.Receive(w.Conn, &msg)
//
//	return fiberws.BinaryMessage, msg, err
//}
//
//// WriteMessage sends a message over the WebSocket connection.
//func (w *webSocketNatsWrapper) WriteMessage(messageType int, data []byte) error {
//	return websocket.Message.Send(w.Conn, data)
//}
//
//// Close the WebSocket connection.
//func (w *webSocketNatsWrapper) Close() error {
//	return w.Conn.Close()
//}
//
//func NewSocketProxy(client *fiberws.Conn, origin string) (*webSocketProxy, *common.AppError) {
//	wsRaw, err := websocket.Dial(config.WsConnectStr, "", origin)
//	if err != nil {
//		return nil, common.ConvertError(err)
//	}
//
//	// struct
//	return &webSocketProxy{
//		// client kết nối vô
//		client: &webSocketClient{IsConnected: false, Conn: client},
//
//		// nats
//		server: &webSocketNatsWrapper{wsRaw},
//	}, nil
//}
//
//func (p *webSocketProxy) Run() {
//	// watting for 5 minutes for token timeout
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5*60)
//	defer func() {
//		p.server.Close()
//		p.client.Close()
//		cancel()
//	}()
//
//	eg, ctx := errgroup.WithContext(ctx)
//
//	eg.Go(func() error {
//		return forwardMessages(ctx, p.client, p.server)
//	})
//
//	eg.Go(func() error {
//		return forwardMessages(ctx, p.server, p.client)
//	})
//
//	var err error
//	go func() {
//		err = eg.Wait()
//	}()
//
//	// tắt kết
//	<-ctx.Done()
//	p.server.Close()
//	p.client.Close()
//
//	if err != nil {
//		fmt.Println("❌ Connection failed", err)
//	}
//
//	fmt.Println("✅ Connection closed")
//}
//
//func forwardMessages(ctx context.Context, src, dst messageTransportIf) error {
//	for {
//		select {
//		case <-ctx.Done():
//			return nil
//		default:
//			msgType, msg, err := src.ReadMessage()
//			fmt.Println(string(msg))
//			if err != nil {
//				return err
//			}
//			if err := dst.WriteMessage(msgType, msg); err != nil {
//				return err
//			}
//		}
//	}
//}
