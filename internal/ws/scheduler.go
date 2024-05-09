package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"websockets/pkg/log"
)

func InitHubScheduler(logger *log.Logs) *HubScheduler {
	return &HubScheduler{
		logger: logger,
		rooms:  &map[string]map[string]*websocket.Conn{},
	}
}

type HubScheduler struct {
	sync.RWMutex
	logger *log.Logs
	rooms  *map[string]map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *HubScheduler) removeClient(roomName, clientUUID string) {
	h.Lock()
	defer h.Unlock()
	delete((*h.rooms)[roomName], clientUUID)
}

func (h *HubScheduler) sendRoomMessages(msgType int, msgBytes []byte, roomName string, senderConn *websocket.Conn) {
	h.RLock()
	defer h.RUnlock()

	for clientUUID, conn := range (*h.rooms)[roomName] {
		if conn == senderConn {
			continue
		}
		err := conn.WriteMessage(msgType, msgBytes)
		if err != nil {
			h.removeClient(roomName, clientUUID)
			h.logger.Error(err.Error())
		}
	}
}

func (h *HubScheduler) listenRoomConnection(roomName, clientUUID string, conn *websocket.Conn) {
	for {
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			h.removeClient(roomName, clientUUID)
			h.logger.Error(err.Error())
			return
		}

		h.sendRoomMessages(msgType, msgBytes, roomName, conn)
	}
}

func (h *HubScheduler) CreateRoom(roomName string, w http.ResponseWriter, r *http.Request) error {
	h.Lock()
	defer h.Unlock()

	if (*h.rooms)[roomName] != nil {
		return RoomAlreadyExists
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	clientUUID := uuid.New().String()

	(*h.rooms)[roomName] = map[string]*websocket.Conn{clientUUID: conn}

	go h.listenRoomConnection(roomName, clientUUID, conn)

	return nil
}

func (h *HubScheduler) JoinRoom(roomName string, w http.ResponseWriter, r *http.Request) error {
	h.Lock()
	defer h.Unlock()

	if (*h.rooms)[roomName] == nil {
		return RoomNotFound
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	clientUUID := uuid.New().String()

	(*h.rooms)[roomName][clientUUID] = conn

	go h.listenRoomConnection(roomName, clientUUID, conn)

	return nil
}
