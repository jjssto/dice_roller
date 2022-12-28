package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func checkDiceArr(diceArr []int8) ([]int8, bool) {
	if len(diceArr) > 200 {
		ret := make([]int8, 200)
		for i := 0; i < MAX_DICE; i++ {
			ret[i] = diceArr[i]
		}
		return ret, false
	} else {
		return diceArr, true
	}
}

func rollDice(c *gin.Context) {
	idStr := c.Param("id")
	var data map[string]string
	err := c.BindJSON(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 32)
	session := sessions.Default(c)
	player := session.Get("player_id").(int)
	col := data["color"]
	char := data["char"]
	mod, _ := strconv.ParseInt(data["mod"], 10, 32)
	action := data["action"]
	dice := data["dice"]
	diceArr := make([]int8, 0)
	if dice != "" {
		json.Unmarshal([]byte(dice), &diceArr)
	}
	diceArr, _ = checkDiceArr(diceArr)
	r := rooms[int(id)]
	r.addPlayer(player, char, col)
	_, err = r.roll(diceArr, int(mod), player, action)
	if err != nil {
		c.Status(http.StatusForbidden)
	} else {
		rooms[int(id)] = r
		c.Status(http.StatusOK)
	}
}

func addPlayer(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	room, _ := strconv.ParseInt(json["room"], 10, 32)
	id, _ := strconv.ParseInt(json["id"], 10, 32)
	name := json["char"]
	color := json["color"]
	r := rooms[int(room)]
	r.addPlayer(int(id), name, color)
	c.Status(http.StatusOK)
}

func addRoomHandler(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	game := json["id"]
	var g Game
	switch game {
	case "CoC":
		g = CoC
	case "RezTech":
		g = RezTech
	default:
		g = General
	}
	id, err := addRoom(g)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.AsciiJSON(http.StatusOK, id)
}

func changeRoomSettings(c *gin.Context) {
	var json map[string]string
	err := c.BindJSON(&json)
	session := sessions.Default(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	playerId := session.Get("player_id").(int)
	roomId, _ := strconv.ParseInt(json["room_id"], 10, 32)
	roomName := json["room_name"]
	color := json["color"]
	room := rooms[int(roomId)]
	if int(playerId) != room.OwnerId {
		c.Status(http.StatusForbidden)
		return
	}
	if color == "-" {
		room.Color = ""
	} else if len(color) > 0 {
		room.Color = color
	}
	if roomName == "-" {
		room.Name = ""
	} else if len(roomName) > 0 {
		room.Name = roomName
	}
	rooms[int(roomId)] = room
	c.Status(http.StatusOK)
}
