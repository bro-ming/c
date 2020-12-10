package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// socket Connection Node
type WsServer struct {
	WsConn    *websocket.Conn //socket connection
	inChan    chan []byte     // receiving chan
	outChan   chan []byte     // release chan
	closeChan chan byte       // close socket

	// 回复ping有效时间
	PingIntervalTime time.Duration

	Mutex   sync.Mutex
	IsClose bool

	// 客户端信息标记 是否登录，身份信息
	Auth bool
	Tag  string
}

// 协议升级器 配置allow cross
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewWsServer 获取通道
func NewWsServer(w http.ResponseWriter, r *http.Request) (conn *WsServer) {
	var (
		wsConn *websocket.Conn
		err    error
	)

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		log.Fatalf("websocket服务机升级协议失败! 原因: %s", err)
		return
	}

	// 本项连接实例
	conn = &WsServer{
		WsConn:    wsConn,
		inChan:    make(chan []byte, 10000),
		outChan:   make(chan []byte, 10000),
		closeChan: make(chan byte, 1),
	}

	// pong 接收及截止设置(可从外部传入,也可读配置文件,暂时写死) 5s内没有ping回复，直接切断
	conn.PingIntervalTime = 3 * time.Second

	// 设置客户端pong回复起始时间，判断是否读取超时
	if err = conn.WsConn.SetReadDeadline(time.Now().Add(conn.PingIntervalTime)); err != nil {
		return
	}

	// 启动读协程
	go conn.readLoop()

	// 启动写协程
	go conn.writeLoop()

	return
}

// Close connection
func (conn *WsServer) Close() {
	_ = conn.WsConn.Close()

	// update WsConn State
	conn.Mutex.Lock()
	if !conn.IsClose {
		close(conn.closeChan)
		conn.IsClose = true
	}
	conn.Mutex.Unlock()
}

// RefreshLife Refresh conn life when client send pongMessage
func (conn *WsServer) RefreshLife() {

	// 接收到pong消息,截止时间向后加
	_ = conn.WsConn.SetReadDeadline(time.Now().Add(conn.PingIntervalTime))
}

// WriteMessage write msg to outChan
func (conn *WsServer) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("websocket服务端客户机断开连接!")
	}
	return
}

// ReadMessage Read message
func (conn *WsServer) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("客户机已断开连接!")
	}
	return
}

// Write msg from OutChan
func (conn *WsServer) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:

			// write message
			if err = conn.WsConn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("客户端断开连接! 客户机信息:", conn.WsConn.RemoteAddr())
				goto ERR
			}

		case <-conn.closeChan:
			log.Print(" 客户端断开连接! 客户机信息:")
			goto ERR
		}
	}

ERR:
	conn.Close()
}

// 内部实现
func (conn *WsServer) readLoop() {
	var (
		data []byte
		err  error
	)

	for {

		if _, data, err = conn.WsConn.ReadMessage(); err != nil {
			log.Print("客户端断开连接，客户机信息：", conn.WsConn.RemoteAddr(), err)
			goto ERR
		}

		//阻塞等待messageType
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			log.Print("客户端断开连接， 客户机信息：", conn.WsConn.RemoteAddr())
			goto ERR
		}
	}

ERR:
	conn.Close()
}
