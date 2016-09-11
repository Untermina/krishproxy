package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"net"
	"net/http"
	"strconv"
	"fmt"
	"encoding/base64"
	"time"
)

const (
	maxMessageSize = 4096
	pongWait = 30 * time.Second
)

type (
	GameServer interface {
		Start(host string, port int)
		Shutdown()
	}

	tcpserver struct {
	}

	webserver struct {
	}

)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewGameServer() GameServer {
	return &webserver{}
}

func (s *tcpserver) Start(host string, port int) {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		log.Errorf("Error getting host:port address: %s", err)
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("Error listening: %s", err)
		return
	}
	log.Infof("Accepting connections on %s", addr)
	for {
		if conn, err := listener.AcceptTCP(); err != nil {
			log.Errorf("Error accepting new connection: %s", err)
			return
		} else {
			log.Infof("New connection from: %s", conn.RemoteAddr().String())
			//go player.Spawn(conn)
		}
	}
}

func (s *tcpserver) Shutdown() {

}

func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

func (s *webserver) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		c.SetReadLimit(maxMessageSize)
		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		err = c.WriteMessage(websocket.TextMessage, []byte("Connecting to game..."))
		if err != nil {
			return
		}

		conn, err := net.Dial("tcp", "127.0.0.1:4000")
		if err != nil {
			fmt.Println(err)
			err = c.WriteMessage(websocket.TextMessage, []byte("Sorry! The game is currently unavailable. Please try again later!"))
			return
		}

		//connbuf := bufio.NewReader(conn)
		closech := make(chan struct{})
		mudSend := make(chan string)
		ticker := time.NewTicker(1500 * time.Millisecond)
		defer func() {
			ticker.Stop()
		}()
		go func() {
			p := make([]byte, 4096)
			for {
				n, err := conn.Read(p)
				if err != nil {
					fmt.Println("Closing socket due to error")
					closech <- struct{}{}
					return
				}
				str := string(p[:n])
				data := []byte(str)
				str = base64.StdEncoding.EncodeToString(data)
				mudSend <- str
			}
		}()
		go func() {
			for {
				// Read
				_, msg, err := c.ReadMessage()
				if err != nil {
					log.Infof("WS closing")
					closech <- struct{}{}
					return
				}
				if len(msg) < 4096 && isASCII(string(msg)) {
					conn.Write(msg)
				} else {
					closech <- struct{}{}
					return
				}
			}
		}()
		for {
			select {
			case <- closech:
				conn.Close()
				return
			case str := <- mudSend:
				err = c.WriteMessage(websocket.TextMessage, []byte(str))
				if err != nil {
					fmt.Println("Bad input.")
					closech <- struct{}{}
					return
				}
			case <-ticker.C:
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					fmt.Println("Ping failed.")
					closech <- struct{}{}
					return
				}
				//conn.Write([]byte("\r\n"))
			}
		}
	}
}

func (s *webserver) Start(host string, port int) {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		log.Errorf("Error getting host:port address: %s", err)
		return
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("./public"))
	e.GET("/ws", standard.WrapHandler(http.HandlerFunc(s.wsHandler())))
	e.Run(standard.New(addr.String()))
}

func (s *webserver) Shutdown() {

}
