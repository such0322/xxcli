package client

import (
	"encoding/json"
	"fmt"
	"xiuxian/common/consts"
	"xiuxian/protocol"
)

var FuncMap = map[string]func(*Game, []byte){
	consts.HandlerWorldGetState: GetState,
	consts.OnWorldPlayerInfo:    PlayerInfo,
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
		w.CircleOne.X = pos.X * 100
		w.CircleOne.Y = pos.Z * 100
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
		w.CircleOne.X = pos.X * 100
		w.CircleOne.Y = pos.Z * 100
	}

}

//func EntitiesState(game *Game, data []byte) {
//
//}
