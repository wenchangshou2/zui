package main

import (
	"fmt"
	"github.com/wenchangshou2/zui/pkg/logging"
	"github.com/wenchangshou2/zui/pkg/setting"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
)

type TouchServer struct {
	backend Backend
	secret  string
}

func (touch *TouchServer) Init(setting setting.Touch) error {
	var (
		err error
	)
	touch.secret = setting.Secret
	touch.backend, err = InitWindowsBackend()
	if err != nil {
		return err
	}
	return nil
}
func (touch *TouchServer) Start() {
	fmt.Println("start1111")
	listener, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("err11",err)
		log.Fatal(err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	host := ""
	bindHost, _, err := net.SplitHostPort(setting.TouchSetting.Bind)
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range addr.IP {
		if b != 0 {
			host = bindHost
			break
		}
	}
	if host == "" {
		host = FindDefaultHost()
	}
	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		fmt.Println("wwww")
		var message string
		for {
			if err := websocket.Message.Receive(ws, &message); err != nil {
				return
			}
			fmt.Println("process command", string(message))
			if err := processCommand(message); err != nil {
				log.Printf(fmt.Sprintf("%s backend:%v", "", err))
				return
			}
		}
	}))
	fmt.Println("start1111")
	http.Serve(listener,mux)
}

func processCommand(message string) error {
	return nil
}
func InitTouchServer(touchSetting *setting.Touch) error {
	logging.G_Logger.Info("init touch server 1")
	touchServer := TouchServer{}
	touchServer.Init(*touchSetting)
	go touchServer.Start()
	fmt.Println("init touch server 4")
	return nil

}
