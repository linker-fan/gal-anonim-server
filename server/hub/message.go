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
const JoinPrivateRoomAction = "join-private-room"
const RoomJoinedAction = "room-joined"

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  *Room   `json:"target"`
	Sender  *Client `json:"sender"`
}

func (m *Message) Encode() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	return json
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	m.Sender = &msg.Sender

	return nil
}
