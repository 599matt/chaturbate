package chaturbate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sextech/chaturbate/internal"
	"github.com/sextech/chaturbate/option"
	"github.com/sextech/chaturbate/utils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Chat struct {
	init         bool
	ctx          context.Context
	conn         *websocket.Conn
	opts         *internal.Settings
	host         string
	username     string
	password     string
	roomPassword string
	room         string
	debug        bool

	OnMessage   func(author, message string)
	OnTip       func(author string, amount int)
	OnNotice    func(notices []string)
	OnUserEntry func(user string)
	OnUserLeave func(user string)
	OnMute      func()
}

type Request struct {
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Args     []interface{} `json:"args"`
	Callback interface{}   `json:"callback"`
	Method   string        `json:"method"`
}

type AuthResult int

const (
	AuthUnknown AuthResult = 0
	AuthOK      AuthResult = 1
)

func NewChat(ctx context.Context, room string, opts ...option.ClientOption) (*Chat, error) {
	var o internal.Settings

	for _, opt := range opts {
		opt.Apply(&o)
	}

	chat := Chat{
		debug: false,
		ctx:   ctx,
		room:  room,
		opts:  &o,
	}

	err := chat.getCredentials()

	if err != nil {
		return nil, err
	}

	chat.init = true

	return &chat, nil
}

func (c *Chat) getCredentials() error {
	if !c.opts.NoAuth {
		return errors.New("not implemented")
	}

	res, err := http.Get(fmt.Sprintf("https://chaturbate.com/%s/", c.room))

	if err != nil {
		return fmt.Errorf("failed to load chaturbate page: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to load chaturbate page, status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("failed to read chaturbate body page: %v", err)
	}

	bodyStr := string(body)

	var reWSChatHost = regexp.MustCompile(`(?m)wschat_host: '(.*)',`)
	var reUsername = regexp.MustCompile(`(?m)username: '(.*)',`)
	var rePassword = regexp.MustCompile(`(?m)password: '(.*)',`)
	var reRoomPassword = regexp.MustCompile(`(?m)room_password: '(.*)',`)

	wsChatHost := reWSChatHost.FindAllStringSubmatch(bodyStr, -1)[0][1]
	username := reUsername.FindAllStringSubmatch(bodyStr, -1)[0][1]
	password := rePassword.FindAllStringSubmatch(bodyStr, -1)[0][1]
	roomPassword := reRoomPassword.FindAllStringSubmatch(bodyStr, -1)[0][1]

	if len(wsChatHost) <= 0 {
		return fmt.Errorf("cam '%s' is not connected", c.room)
	}

	wsChatHost = strings.Replace(wsChatHost, "https://", "wss://", 1)

	c.host = wsChatHost
	c.username = username
	c.password = password
	c.roomPassword = roomPassword

	return nil
}

func (c *Chat) Connect() error {
	var err error

	if !c.init {
		return errors.New("chat client not initialized")
	}

	url := fmt.Sprintf("%s/%d/%s/websocket", c.host, rand.Intn(1000), utils.RandomString(8))

	c.conn, _, err = websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return fmt.Errorf("failed to dial websocket server: %v", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := c.conn.ReadMessage()

			if err != nil {
				log.Println("failed to read:", err)
				return
			}

			err = c.parse(message)

			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil
		}
	}
}

func (c *Chat) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return errors.New("chat is not connected")
}

func (c *Chat) parse(message []byte) error {
	if len(message) <= 0 {
		return errors.New("message is empty")
	}

	// first letter is command
	switch string(message[0]) {
	case "o":
		return c.onConnected()
	case "h":
		// No response needed to 'h' message
	case "a":
		var data []string

		// skip first letter
		err := json.Unmarshal(message[1:], &data)

		if err != nil {
			return fmt.Errorf("failed to parse raw data: %v", err)
		}

		if len(data) <= 0 {
			return errors.New("no data")
		}

		var resp Response

		err = json.Unmarshal([]byte(data[0]), &resp)

		if err != nil {
			return fmt.Errorf("failed to parse data: %v", err)
		}

		switch resp.Method {
		case "onAuthResponse":
			resultStr := resp.Args[0].(string)
			result, _ := strconv.Atoi(resultStr)
			return c.onAuthResponse(AuthResult(result))
		case "onNotify":
			dataStr := resp.Args[0].(string)

			var notify notify
			err := json.Unmarshal([]byte(dataStr), &notify)

			if err != nil {

			}

			return c.onNotify(notify)
		case "onTitleChange":
			// TODO
		case "onRoomCountUpdate":
			// TODO
		case "onRoomMsg":
			author := resp.Args[0].(string)
			dataStr := resp.Args[1].(string)

			var roomMessage roomMessage
			err := json.Unmarshal([]byte(dataStr), &roomMessage)

			if err != nil {
				return fmt.Errorf("failed to parse onRoomMsg data: %v", err)
			}

			return c.onRoomMsg(author, roomMessage)
		case "onSilence":
			// TODO
		default:
			fmt.Printf("unknown message: %s -> %s\n", resp.Method, message)
		}

	default:
		log.Printf("unknown message type: %s\n", message)
	}

	return nil
}

func (c *Chat) send(method string, data interface{}) error {
	req := Request{
		Method: method,
		Data:   data,
	}

	requestJson, err := json.Marshal(req)

	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	var array []string
	array = append(array, string(requestJson))

	toSend, err := json.Marshal(array)

	if err != nil {
		return fmt.Errorf("failed to marshal request array: %v", err)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, toSend)

	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (c *Chat) onConnected() error {
	connectMessage := connectMessage{
		User:         c.username,
		Password:     c.password,
		Room:         c.room,
		RoomPassword: c.roomPassword,
	}

	return c.send("connect", connectMessage)
}

func (c *Chat) onAuthResponse(result AuthResult) error {
	if result == AuthOK {
		joinRoomMessage := joinRoomMessage{
			Room: c.room,
		}

		return c.send("joinRoom", joinRoomMessage)
		// TODO call OnConnected callback
	}

	return fmt.Errorf("bad auth response result: %d", result)
}

func (c *Chat) onNotify(notify notify) error {
	switch notify.Type {
	case "room_entry":
		if c.OnUserEntry != nil {
			go c.OnUserEntry(notify.Username)
		}
	case "room_leave":
		if c.OnUserLeave != nil {
			go c.OnUserLeave(notify.Username)
		}
	case "purchase_notification":
		// TODO
	case "appnotice":
		if c.OnNotice != nil {
			go c.OnNotice(notify.Msg)
		}
	case "refresh_panel":
		// TODO
	case "tip_alert":
		if c.OnTip != nil {
			go c.OnTip(notify.FromUsername, notify.Amount)
		}
	default:
		log.Printf("unknown notification type: %s -> %+v\n", notify.Type, notify)
	}

	return nil
}

func (c *Chat) onRoomMsg(author string, roomMessage roomMessage) error {
	if c.OnMessage != nil {
		go c.OnMessage(author, roomMessage.Message)
	}

	return nil
}
