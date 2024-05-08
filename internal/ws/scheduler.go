package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"websockets/pkg/log"
)

func InitHubScheduler(logger *log.Logs) *HubScheduler {
	return &HubScheduler{
		logger: logger,
		rooms:  &map[string][]*websocket.Conn{},
	}
}

type HubScheduler struct {
	logger *log.Logs
	rooms  *map[string][]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *HubScheduler) sendRoomMessages(msgType int, msgBytes []byte, roomName string, senderConn *websocket.Conn) {
	for _, conn := range (*h.rooms)[roomName] {
		if conn == senderConn {
			continue
		}
		err := conn.WriteMessage(msgType, msgBytes)
		if err != nil {
			h.logger.Error(err.Error())
		}
	}
}

func (h *HubScheduler) listenRoomConnection(roomName string, conn *websocket.Conn) {
	for {
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error(err.Error())
			return
		}

		h.sendRoomMessages(msgType, msgBytes, roomName, conn)
	}
}

func (h *HubScheduler) CreateRoom(roomName string, w http.ResponseWriter, r *http.Request) error {
	if (*h.rooms)[roomName] != nil {
		return RoomAlreadyExists
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	(*h.rooms)[roomName] = []*websocket.Conn{conn}

	go h.listenRoomConnection(roomName, conn)

	return nil
}

func (h *HubScheduler) JoinRoom(roomName string, w http.ResponseWriter, r *http.Request) error {
	if (*h.rooms)[roomName] == nil {
		return RoomNotFound
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	(*h.rooms)[roomName] = append((*h.rooms)[roomName], conn)

	go h.listenRoomConnection(roomName, conn)

	return nil
}
