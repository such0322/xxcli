package client

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/v2/client"
	"github.com/topfreegames/pitaya/v2/session"
	"time"
)

const (
	SendIntervalTime = 50
)

const (
	SendMsgTypReq = 1
	SendMsgTypNtf = 2
)

type sendMsg struct {
	typ   int
	route string
	data  proto.Message
}

func NewSendMsg(route string, data proto.Message, typ int) *sendMsg {
	return &sendMsg{
		typ:   typ,
		route: route,
		data:  data,
	}
}

var agent *Agent

type Agent struct {
	Cli     client.PitayaClient
	MsgMap  map[uint]string
	SendMsg chan *sendMsg
}

func NewAgent() *Agent {
	a := &Agent{}
	a.Cli = client.New(logrus.InfoLevel)
	a.MsgMap = make(map[uint]string)
	a.SendMsg = make(chan *sendMsg, 128)
	return a
}

func InitAgent() {
	agent = NewAgent()
	handshake := &session.HandshakeData{
		Sys: session.HandshakeClientData{
			Platform:    "win",
			LibVersion:  "0.3.5-release",
			BuildNumber: "20",
			Version:     "1.0.0",
		},
		User: map[string]interface{}{
			"age": 30,
		},
	}
	agent.Cli.SetClientHandshakeData(handshake)
	go SendServerMessages()
}

func Disconnect() {
	if agent.Cli.ConnectedStatus() {
		fmt.Println("disconnect")
		agent.Cli.Disconnect()
	}
}

func SendServerMessages() {
	t := time.NewTicker(SendIntervalTime * time.Millisecond)
	ntfMsg := make(map[string]proto.Message)
	reqMsg := make(map[string]proto.Message)
	for {
		select {
		case <-t.C:
			for r, req := range reqMsg {
				agent.request(r, req)
				delete(reqMsg, r)
			}
			for n, ntf := range ntfMsg {
				agent.notify(n, ntf)
				delete(ntfMsg, n)
			}
		case msg := <-agent.SendMsg:
			if msg.typ == SendMsgTypReq {
				reqMsg[msg.route] = msg.data
			}
			if msg.typ == SendMsgTypNtf {
				ntfMsg[msg.route] = msg.data
			}
		}

	}
}

func ReadServerMessages(game *Game) {
	channel := agent.Cli.MsgChannel()
	for {
		select {
		case m := <-channel:
			if m.Err {
				fmt.Println("read server msg is err msg, err:", m.Err, " msgId =", m.ID, " route =", m.Route, " data =", string(m.Data))
				continue
			}
			if m.Type == 2 {
				m.Route = agent.MsgMap[m.ID]
			}
			fmt.Println(m.Type, m.ID, m.Route)
			fmt.Println("Data =", string(m.Data))
			fn, ok := FuncMap[m.Route]
			if !ok {
				fmt.Println("route func =", m.Route, "not found")
			} else {
				fn(game, m.Data)
			}
			fmt.Println("---------------------------------------------------------------")
		}
	}
}

func Request(route string, data proto.Message) {
	msg := NewSendMsg(route, data, SendMsgTypReq)
	agent.SendMsg <- msg
}

func (*Agent) request(route string, data proto.Message) {
	b, _ := json.Marshal(data)
	mid, err := agent.Cli.SendRequest(route, b)
	if err != nil {
		fmt.Println("agent request err: ", err)
	}
	agent.MsgMap[mid] = route
}

func Notify(route string, data proto.Message) {
	msg := NewSendMsg(route, data, SendMsgTypNtf)
	agent.SendMsg <- msg

}

func (*Agent) notify(route string, data proto.Message) {
	b, _ := json.Marshal(data)
	err := agent.Cli.SendNotify(route, b)
	if err != nil {
		fmt.Println("agent notify err: ", err)
	}
}
