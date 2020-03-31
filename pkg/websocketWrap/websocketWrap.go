package websocketWrap

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
)

var ErrNotConnected = errors.New("websocket: not connected")

type RecConn struct {
	//指定重连的间隔
	// 默认2秒
	RecIntvlMin time.Duration
	//指定最大重新连接间隔
	// 默认30秒
	RecIntvlMax time.Duration
	//指定重连的增加量
	//默认1.5
	RecIntvlFactor float64
	// HandshakeTimeout指定握手完成的持续时间，
	HandshakeTimeout time.Duration
	NonVerbose       bool
	SubscribeHandler func() error
	KeepAliveTimeout time.Duration

	mu          sync.RWMutex
	url         string
	reqHeader   http.Header
	httpResp    *http.Response
	dialErr     error
	isConnected bool
	dialer      *websocket.Dialer
	regFunc     func() bool

	*websocket.Conn
}

func (rc *RecConn) closeAndReconnect() {
	rc.Close()
	go rc.connect()
}
func (rc *RecConn) getConn() *websocket.Conn {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.Conn
}

//func (rc *RecConn) SetConnectSuccessFunc(f func() bool){
//	rc.regFunc=f
//}
func (rc *RecConn) Close() {
	if rc.getConn() != nil {
		rc.mu.Lock()
		rc.Conn.Close()
		rc.mu.Unlock()
	}
	rc.setIsConnected(false)
}
func (rc *RecConn) setIsConnected(state bool) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.isConnected = state
}

func (rc *RecConn) getBackoff() *backoff.Backoff {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return &backoff.Backoff{
		Min:    rc.RecIntvlMin,
		Max:    rc.RecIntvlMax,
		Factor: rc.RecIntvlFactor,
		Jitter: true,
	}
}

func (rc *RecConn) ReadMessage() (messageType int, message []byte, err error) {
	err = ErrNotConnected
	if rc.IsConnected() {
		messageType, message, err = rc.Conn.ReadMessage()
		if err != nil {
			rc.closeAndReconnect()
		}
	}
	return
}
func (rc *RecConn) readJSON(v interface{}) error {
	err := ErrNotConnected
	if rc.IsConnected() {
		_, message, err := rc.Conn.ReadMessage()
		if err != nil {
			rc.closeAndReconnect()
		} else {
			json.Unmarshal(message, &v)
		}
	}
	return err
}
func (rc *RecConn) getNonVerbose() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return rc.NonVerbose
}

func (rc *RecConn) hasSubscribeHandler() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return rc.SubscribeHandler != nil
}
func (rc *RecConn) getKeepAliveTimeout() time.Duration {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.KeepAliveTimeout
}
func (rc *RecConn) parseURL(urlStr string) (string, error) {
	if urlStr == "" {
		return "", errors.New("dial: url cannot be empty")
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", errors.New("url: " + err.Error())
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return "", errors.New("url:websocket uris must start with ws or wss scheme")
	}
	if u.User != nil {
		return "", errors.New("url: user name and password are not allowed in websocket URIs")
	}
	return urlStr, nil
}
func (rc *RecConn) writeControlPingMessage() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	return rc.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
}
func (rc *RecConn) SetWriteDeadline(t time.Duration) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.isConnected {
		return rc.Conn.SetWriteDeadline(time.Now().Add(t))
	}
	return nil

}
func (rc *RecConn) keepAlive() {
	var (
		keepAliveResponse = new(keepAliveResponse)
		ticker            = time.NewTicker(rc.getKeepAliveTimeout())
	)

	rc.mu.Lock()
	rc.Conn.SetPongHandler(func(msg string) error {
		keepAliveResponse.setLastResponse()
		return nil
	})
	rc.mu.Unlock()

	go func() {
		defer ticker.Stop()

		for {
			rc.writeControlPingMessage()
			<-ticker.C
			if time.Now().Sub(keepAliveResponse.getLastResponse()) > rc.getKeepAliveTimeout() {
				rc.closeAndReconnect()
				return
			}
		}
	}()
}

func (rc *RecConn) setURL(url string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.url = url
}
func (rc *RecConn) setReqHeader(reqHeader http.Header) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.reqHeader = reqHeader
}
func (rc *RecConn) setDefaultRecIntvlMin() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.RecIntvlMin == 0 {
		rc.RecIntvlMin = 2 * time.Second
	}
}

func (rc *RecConn) setDefaultRecIntvlMax() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.RecIntvlMax == 0 {
		rc.RecIntvlMax = 30 * time.Second
	}
}
func (rc *RecConn) setDefaultRecIntvlFactor() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.RecIntvlFactor == 0 {
		rc.RecIntvlFactor = 1.5
	}
}

func (rc *RecConn) setDefaultHandshakeTimeout() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.HandshakeTimeout == 0 {
		rc.HandshakeTimeout = 2 * time.Second
	}
}
func (rc *RecConn) setDefaultDialer(handshakeTimeout time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.dialer = &websocket.Dialer{
		HandshakeTimeout: handshakeTimeout,
	}
}
func (rc *RecConn) getHandshakeTimeout() time.Duration {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.HandshakeTimeout
}
func (rc *RecConn) Dial(urlStr string, reqHeader http.Header) {
	urlStr, err := rc.parseURL(urlStr)
	if err != nil {
		log.Fatalf("Dial: %v", err)
	}

	// Config
	rc.setURL(urlStr)
	rc.setReqHeader(reqHeader)
	rc.setDefaultRecIntvlMin()
	rc.setDefaultRecIntvlMax()
	rc.setDefaultRecIntvlFactor()
	rc.setDefaultHandshakeTimeout()
	rc.setDefaultDialer(rc.getHandshakeTimeout())

	// Connect
	go rc.connect()

	// wait on first attempt
	time.Sleep(rc.getHandshakeTimeout())
}
func (rc *RecConn) connect() {
	b := rc.getBackoff()
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		nextItvl := b.Duration()
		wsConn, httpResp, err := rc.dialer.Dial(rc.url, rc.reqHeader)

		rc.mu.Lock()
		rc.Conn = wsConn
		rc.dialErr = err
		rc.isConnected = err == nil
		rc.httpResp = httpResp
		rc.mu.Unlock()
		if err == nil {
			if !rc.getNonVerbose() {
				log.Printf("Dial: connection was successfully established with %s\n", rc.url)

				if !rc.hasSubscribeHandler() {
					return
				}

				if err := rc.SubscribeHandler(); err != nil {
					log.Fatalf("Dial: connect handler failed with %s", err.Error())
				}

				log.Printf("Dial: connect handler was successfully established with %s\n", rc.url)

				if rc.getKeepAliveTimeout() != 0 {
					rc.keepAlive()
				}
			}

			return
		}
		if !rc.getNonVerbose() {
			log.Println(err)
			log.Println("Dial: will try again in", nextItvl, "seconds.")
		}
		time.Sleep(nextItvl)
	}
}
func (rc *RecConn) IsConnected() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.isConnected
}
