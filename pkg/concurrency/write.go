package main

//
//import (
//	"github.com/gorilla/websocket"
//	"time"
//)
//
//func (c *conn) handleWrite() {
//	ticker := time.NewTicker(_pingPeriod)
//	defer func() {
//		ticker.Stop()
//		c.Close(xcontext.NewTraceContext(context.TODO(), "close_ws_connection"))
//	}()
//
//	for {
//		select {
//		case message, ok := <-c.outgoingMessages:
//			c.resetWriteDeadline()
//			if !ok {
//				c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//			w, err := c.wsConn.NextWriter(message.msgType)
//			if err != nil {
//				logging.Logger(message.ctx).Error("next writer failed",
//					zap.String("conn_id", c.connID),
//					zap.String("safe_id", c.safeID),
//					zap.Error(err))
//				return
//			}
//			if _, err = w.Write(message.payload); err != nil {
//				logging.Logger(message.ctx).Error("write message failed",
//					zap.String("conn_id", c.connID),
//					zap.Error(err))
//			}
//			if err := w.Close(); err != nil {
//				logging.Logger(message.ctx).Error("close writer failed",
//					zap.String("conn_id", c.connID),
//					zap.String("safe_id", c.safeID),
//					zap.Error(err))
//				return
//			}
//		case <-ticker.C:
//			c.wsConn.SetWriteDeadline(time.Now().Add(_writeDeadline))
//			if err := c.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				return
//			}
//		}
//	}
//}
