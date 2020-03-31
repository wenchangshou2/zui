package main

/*
	#cgo windows LDFLAGS: -lgdi32 -luser32
	#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/cdeps/win64 -lpng -lz
	#cgo windows,386 LDFLAGS: -L${SRCDIR}/cdeps/win32 -lpng -lz
	#include "window/Window.h"
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/wenchangshou2/zoolon_message"
	"github.com/wenchangshou2/zui/form"
	"github.com/wenchangshou2/zui/pkg/logging"
	"github.com/wenchangshou2/zui/pkg/websocketWrap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024 * 1
)

//一个连接
type Client struct {
	Ip            string
	Port          int
	connectStatus bool
	conn          *websocket.Conn
	send          chan []byte
	memoryMsgChan chan *zoolon_message.Message
	websocketWrap.RecConn
	lastReportLayoutInfoTime int64
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	validate *validator.Validate
	G_Client *Client
)


func (c *Client) active(data []byte) error {
	//var (
	//	cmd e.ActiveWindowForm
	//	err error
	//)
	//logging.G_Logger.Debug(fmt.Sprintf("active function"))
	//cmd=e.ActiveWindowForm{}
	//err=json.Unmarshal(data,&cmd)
	//if err!=nil{
	//	logging.G_Logger.Error(fmt.Sprintf("active window parse json failed:%s",err.Error()))
	//	return err
	//}
	//pid:=G_Vm.GetPidByWid(cmd.Wid)
	//if pid==-1{
	//	logging.G_Logger.Info("当前窗口的线程未找到")
	//	return errors.New("当前窗口的线程未找到")
	//}
	//exists,err:=process.PidExists(pid)
	//if !exists||err!=nil{
	//	logging.G_Logger.Info("窗口进程异常")
	//	return errors.New("窗口的进程异常")
	//}
	//logging.G_Logger.Info("active window success")
	//return robotgo.ActivePID(pid)
	return nil
}
func (c *Client) execute(action string, SenderName string, data []byte) (map[string]interface{}, error) {
	validate = validator.New()
	var (
		err     error
		rtuData map[string]interface{}
	)
	logging.G_Logger.Info(fmt.Sprintf("action:%s", action))
	switch action {
	case "openMultiScreen":
		fmt.Println("")
	default:
		err = errors.New("未支持的操作")
	}
	rtu := make(map[string]interface{})
	rtu["state"] = 200
	rtu["message"] = "成功"
	if err != nil {
		rtu["state"] = 400
		rtu["message"] = err.Error()
	}
	rtu["receiverName"] = "/daemon"
	rtu["senderName"] = "/zui"
	rtu["Action"] = action
	if len(rtuData) > 0 {
		rtu["data"] = rtuData
	}
	return rtu, nil
}
func (c *Client) binaryExecute(action string, SenderName string, data []byte) (map[string]interface{}, string, error) {
	validate = validator.New()
	var (
		err     error
		rtuData interface{}
	)
	topic := "/zebus"
	rtu := make(map[string]interface{})
	logging.G_Logger.Info(fmt.Sprintf("action:%s", action))
	switch action {
	case "move":
		err = c.move(data)
	case "scroll":
		err = c.scroll(data)
	case "active":
		err = c.activeByPid(data)
	case "keyboard":
		err = c.keyboard(data)
	case "queryLayout":
	default:
		err = errors.New("未支持的操作111")
	}
	rtu["state"] = 200
	rtu["message"] = "成功"
	if err != nil {
		rtu["state"] = 400
		rtu["message"] = err.Error()
	}
	rtu["Action"] = action
	if rtuData != nil {
		rtu["data"] = rtuData
	}
	return rtu, topic, err
}

// 二进制消息处理
func (c *Client) BinaryMessageProcess(message *zoolon_message.Message) (err error) {
	jsonData := map[string]interface{}{}
	err = json.Unmarshal(message.Body, &jsonData)
	if err != nil {
		logging.G_Logger.Warn(fmt.Sprintf("parse json error:(%s)", err.Error()))
		return err
	}
	if action, ok := jsonData["Action"]; ok {
		var senderName = ""
		if _, ok = jsonData["senderName"]; ok {
			senderName = jsonData["senderName"].(string)
		}
		result, topic, err := c.binaryExecute(action.(string), senderName, message.Body)
		fmt.Println("result", result, topic, err)
		json, _ := json.Marshal(result)
		logging.G_Logger.Info("send message:" + string(json))
		message.SetBody(string(json))
		message.SetTopic(topic)
		c.memoryMsgChan <- message
	}
	return
}

// 文本消息处理
func (c *Client) TextMessageProcess(message []byte) (err error) {
	jsonData := map[string]interface{}{}
	err = json.Unmarshal(message, &jsonData)
	if err != nil {
		fmt.Errorf("解析json失败")
	}
	if value, ok := jsonData["Action"]; ok {
		var senderName = ""
		if _, ok = jsonData["senderName"]; ok {
			senderName = jsonData["senderName"].(string)
		}
		c.execute(value.(string), senderName, message)
	}
	return
}
func (c *Client) ReadPump() {

	for {
		if !c.IsConnected() {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		messageType, message, err := c.ReadMessage()
		if err != nil {
			logging.G_Logger.Warn(fmt.Sprintf("recv message field:(%s)", err.Error()))
			continue
		}
		if messageType == websocket.BinaryMessage {
			msg, err := zoolon_message.DecodeMessage(message)
			if err != nil {
				logging.G_Logger.Info(fmt.Sprintf("解析message 失败:%v", err))
				continue
			}
			logging.G_Logger.Info("recv message:" + string(msg.Body))
			err = c.BinaryMessageProcess(msg)
		} else {
			c.TextMessageProcess(message)

		}
	}
}

func (c *Client) Write(msg interface{}) {
	logging.G_Logger.Info(fmt.Sprintf("write msg:%v", msg))
	c.WriteJSON(msg)
}
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if c.IsConnected() {
			c.Close()
		}
	}()
	for {
		if !c.IsConnected() {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		select {
		case message, ok := <-c.send:
			logging.G_Logger.Info(fmt.Sprintf("write msg success:%s,%b", string(message), ok))
			c.SetWriteDeadline(writeWait)
			if !ok {
				c.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				logging.G_Logger.Info(fmt.Sprintf("websocket close:%v", err))
				return
			}
		case msg, ok := <-c.memoryMsgChan:
			if !c.IsConnected() {
				continue
			}
			var buf = &bytes.Buffer{}
			_, err := msg.WriteTo(buf)
			logging.G_Logger.Info("write message:" + string(msg.Body) + ";write topic:" + string(msg.Topic))
			if err != nil {
				logging.G_Logger.Error(fmt.Sprintf("解析memoryMsgChan错误:%v", err))
				continue
			}
			c.SetWriteDeadline(writeWait)
			if !ok {
				c.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(buf.Bytes())
			w.Write(newline)
			if err := w.Close(); err != nil {
				logging.G_Logger.Info("send memoryMsg error:" + err.Error())
				return
			}
			logging.G_Logger.Info("send memoryMsg success")
		}
	}
}

//
func (c *Client) Loop() {
	for {
		time.Sleep(1 * time.Second)
	}
}
func (c *Client) Start() (err error) {
	//c.SetConnectSuccessFunc(c.RegFunc)
	c.SubscribeHandler = func() error {
		//{"messageType":"RegisterToDaemon","SocketName":"自己程序名称，需要各自不同","SocketType":"Controller"}
		data := map[string]interface{}{}
		data["messageType"] = "RegisterToDaemon"
		data["SocketName"] = "zui"
		data["SocketType"] = "Server"
		data["proto"] = "binary"
		c.WriteJSON(data)
		return nil
	}
	c.Dial(fmt.Sprintf("ws://%s:%d", c.Ip, c.Port), nil)

	go c.ReadPump()
	go c.writePump()
	go c.Loop()
	return
}

func (c *Client) move(data []byte) error {
	var (
		r form.MoveMouseStruct
	)
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	//G_Backend.PointerMove(r.X,r.Y)
	robotgo.Move(r.X, r.Y)
	return nil
}

func (c *Client) scroll(data []byte) error {
	var (
		r form.ScrollMouseRequestBody
	)
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	G_Backend.PointerScroll(r.Horizontal, r.Vertical)
	return nil

}

func (c *Client) activeByPid(data []byte) error {
	var (
		r form.ActiveWindowByPidRequestBody
	)
	if err := json.Unmarshal(data, &r); err != nil {
		logging.G_Logger.Info("解析json 失败:" + string(data))
		return err
	}
	if r.Data.Pid <= 0 {
		logging.G_Logger.Info("pid 必须存在")
		return errors.New("pid 必须存在")
	}
	C.active_force_pid()

	logging.G_Logger.Info(fmt.Sprintf("active wid :%d", r.Data.Pid))
	//processList, err := ps.Processes()
	//logging.G_Logger.Info(fmt.Sprintf("activepid:%d",r.Data.Pid))
	//robotgo.ActivePID(r.Data.Pid)
	robotgo.ActivePID(r.Data.Pid)

	handle := robotgo.GetHandle()
	logging.G_Logger.Info(fmt.Sprintf("handle:%s", handle))
	mdata := robotgo.GetActive()
	robotgo.SetActive(mdata)
	return nil
}

func (c *Client) keyboard(data []byte) error {
	var (
		r form.KeyboardRequestBody
	)
	if err := json.Unmarshal(data, &r); err != nil {
		logging.G_Logger.Info("解析json 失败:" + string(data))
		return err
	}
	if len(r.Key) == 0 {
		return errors.New("key 必须存在")
	}
	keyArr := strings.Split(r.Key, "+")
	if len(keyArr) == 1 {
		robotgo.KeyTap(keyArr[0])
	} else {

		robotgo.KeyTap(keyArr[0], keyArr[1:])
	}
	fmt.Println(r)
	return nil
}

//初始化调度
func InitSchedule(Ip string, port int) (err error) {
	G_Client = &Client{
		Ip:            Ip,
		Port:          port,
		memoryMsgChan: make(chan *zoolon_message.Message, 0),
	}
	go G_Client.Start()
	return
}
