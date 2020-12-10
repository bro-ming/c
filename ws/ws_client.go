package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type WsClient struct {
	sync.Mutex
	WsName      string
	WsUrl       string
	Proxy       string
	IsClose     bool
	CloseChan   chan byte
	PingHandler func(string) error
	Conn        *websocket.Conn // WS客户端

	inChan  chan []byte // receiving chan
	outChan chan []byte //  release chan

	PongChan   chan []byte
	Reconn     int64 // 三方连接信号 0:首次连接,无异常 1:断开 2:重连成功
	ReconnFail bool  // 重连失败
	Login      bool  //是否已登录
}

// NewWsClient 初始化websocket客户端
func NewWsClient(wsUrl string) (ws *WsClient, err error) {
	ws = &WsClient{
		WsUrl:     wsUrl,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		PongChan:  make(chan []byte, 1),
		CloseChan: make(chan byte, 1),
	}

	if err = ws.connect(); err != nil {
		log.Fatalf("websocekt创建失败%s", err)
		return
	}

	// ping 消息接受及次数处理
	ws.Conn.SetPingHandler(func(ping string) error {
		ws.PongChan <- []byte(ping)
		return nil
	})

	go ws.readLoop()
	go ws.writeLoop()
	return
}

// connect 建立websocket连接
func (this *WsClient) connect() error {
	this.Lock()
	defer this.Unlock()

	// 获得拨号实例
	dialer := websocket.DefaultDialer

	// 如果已配置代理播报，从代理拨出去
	if this.Proxy != "" {
		proUrl, _ := url.Parse(this.Proxy)
		dialer.Proxy = http.ProxyURL(proUrl)
	}

	// 发起拨号
	conn, _, err := dialer.Dial(this.WsUrl, nil)
	if err != nil {
		return err
	}

	// 重连以后,更新wsConn,部分属性初始化
	this.Conn = conn
	return nil
}

// CloseWs 主动关闭连接
func (this *WsClient) CloseWs() {

	// 关闭socekt连接
	_ = this.Conn.Close()

	// 关闭已开启的关闭通道
	this.Lock()
	if !this.IsClose {
		this.IsClose = true

		// 断开
		this.Reconn = 1
		close(this.CloseChan)
	}
	this.Unlock()
}

// Reconnect 断开重连
func (this *WsClient) Reconnect() {
	var err error

	// Reconnect 断开重连
	for retry := 1; retry <= 4; retry++ {
		if err = this.connect(); err != nil {
			log.Printf(this.WsName+"重连失败, %s", err.Error())
		} else {

			this.IsClose = false
			this.CloseChan = make(chan byte, 1)
			this.ReconnFail = false
			this.Reconn = 2
			log.Printf(this.WsName + "重连成功!")
			break
		}

		time.Sleep(time.Second * time.Duration(retry))
	}

	if err != nil {
		this.ReconnFail = true
		log.Printf(this.WsName + "多次重连失败,关闭连接!. ")
	}
}

// WriteMessage write msg to outChan
func (conn *WsClient) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.CloseChan:
		err = errors.New("websocket服务端客户机断开连接!")
	}
	return
}

// ReadMessage Read message
func (conn *WsClient) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.CloseChan:
		err = errors.New("服务端客户机已断开连接,请重新连接!")
	}
	return
}

// Write msg from OutChan
func (conn *WsClient) writeLoop() {
	var (
		data []byte
		err  error
	)

	// write message
	for {
		select {
		case data = <-conn.outChan:
			if err = conn.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("websocket服务端客户机断开连接! 客户机信息:%s", conn.WsName)

				// 发起重连
				conn.Reconnect()
				continue
			}

		case data = <-conn.PongChan:
			if err = conn.Conn.WriteMessage(websocket.PongMessage, data); err != nil {
				log.Printf("websocket服务端客户机断开连接! 客户机信息:%s", conn.WsName)

				// 发起重连
				conn.Reconnect()
				continue
			}

		case <-conn.CloseChan:
			log.Printf("websocket服务端客户机断开连接! 客户机信息:%s", conn.WsName)
			goto ERR
		}
	}

ERR:
	conn.CloseWs()
}

// 内部实现
func (conn *WsClient) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.Conn.ReadMessage(); err != nil {
			log.Printf("服务机断开,信息:%s - %s - %s", conn.WsName, conn.Conn.RemoteAddr(), err)

			// 发起重连
			conn.Reconnect()
			continue
		}

		//阻塞等待messageType
		select {
		case conn.inChan <- data:
		case <-conn.CloseChan:
			log.Printf("websocket服务端客户机断开连接! 客户机信息:%s", conn.Conn.RemoteAddr())
			goto ERR
		}
	}

ERR:
	conn.CloseWs()
}
