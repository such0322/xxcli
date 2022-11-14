package client

import (
	"encoding/json"
	"fmt"
	"xiuxian/protocol"
)

const (
	GateLogin        = "gate.entryhandler.login"
	PublicSelectRole = "public.playerhandler.selectrole"
	GameStartGame    = "game.playerhandler.startgame"
	WorldMoveToNtf   = "world.playerhandler.moveto"
	WorldGetState    = "world.playerhandler.getstate"
)

var FuncMap = map[string]func(*Game, []byte){
	"world.playerhandler.getstate": GetState,
	"onworldplayerinfo":            PlayerInfo,
	//"onworldentitiesstate":         EntitiesState,
}

func GetState(game *Game, data []byte) {
	msg := &protocol.MsgWorldGetStateRsp{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println("json unmarshal decode err", err)
		return
	}
	pos := msg.Player.Pos
	w := game.World.(*World)
	if w.Run {
		w.CircleOne.X = pos.X
		w.CircleOne.Y = pos.Z
	}
}

func PlayerInfo(game *Game, data []byte) {
	msg := &protocol.MsgWorldPlayerInfoPush{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println("json unmarshal decode err", err)
		return
	}
	fmt.Println("player pos change")
	pos := msg.Player.Pos
	w := game.World.(*World)
	if w.Run {
		w.CircleOne.X = pos.X
		w.CircleOne.Y = pos.Z
	}

}

//func EntitiesState(game *Game, data []byte) {
//
//}
