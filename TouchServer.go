package main

import (
	"errors"
	"fmt"
	"github.com/wenchangshou2/zui/pkg/computer"
	"github.com/wenchangshou2/zui/pkg/logging"
	"github.com/wenchangshou2/zui/pkg/setting"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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
	listener, err := net.Listen("tcp", "0.0.0.0:8889")
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
	http.Serve(listener,mux)
}

func processCommand(command string) error {
	if len(command)==0{
		return errors.New("empty command")
	}
	arguments := strings.Split(command[1:], ";")
	if command == "sf"{
		return G_Backend.PointerScrollFinish()
	}
	if command[0]=='t'{
		text:=command[1:]
		text=strings.Replace(text,"\r\n","\n",-1)
		text=strings.Replace(text,"\r","\n",-1)
		if!utf8.ValidString(text){
			return errors.New("invalid utf-8")
		}
		return G_Backend.KeyboardText(text)
	}
	if command[0]=='k'{
		fmt.Println("len",len(arguments),arguments)
		argLen:=len(arguments)
		if argLen==0{
			return errors.New("键值不存在")
		}
		if argLen==1{
			fmt.Println("111",arguments[0])
			G_Backend.KeyTap(arguments[0])
		}else if argLen==2{
			G_Backend.KeyTap(arguments[0],arguments[1])
		}else if argLen==3{
			G_Backend.KeyTap(arguments[0],arguments[1],arguments[2])
		}
		return nil
	}
	x, err := strconv.ParseInt(arguments[0], 10, 32)
	if err != nil {
		return err
	}
	y, err := strconv.ParseInt(arguments[1], 10, 32)
	if err != nil {
		return err
	}
	if command[0]=='m'{
		return G_Backend.PointerMove(int(x),int(y))
	}
	if command[0]=='b'{
		if x<0||x>=int64(PointerButtonLimit){
			return errors.New("unsupported pointer button")
		}
		b:=true
		if y==0{
			b=false
		}
		return G_Backend.PointerButton(computer.PointerButton(x),b)
	}
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
