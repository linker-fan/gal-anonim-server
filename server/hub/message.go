package hub

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  string  `json:"target"`
	Sender  *Client `json:"sender"`
}

func (m *Message) encode() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	return json
}
