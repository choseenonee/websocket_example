package scheduler

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
	"websockets/internal/models"
	"websockets/internal/repository"
	"websockets/pkg/log"
)

func InitHubScheduler(logger *log.Logs, repoMessageCreator RepoMessageCreator,
	repoChatGetter repository.ChatGetterRepo) *HubScheduler {
	return &HubScheduler{
		logger:             logger,
		chats:              &map[int]map[string]*websocket.Conn{},
		repoMessageCreator: repoMessageCreator,
		repoChatGetter:     repoChatGetter,
	}
}

type HubScheduler struct {
	sync.RWMutex
	logger             *log.Logs
	chats              *map[int]map[string]*websocket.Conn
	repoChatGetter     repository.ChatGetterRepo
	repoMessageCreator RepoMessageCreator
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *HubScheduler) removeClient(chatID int, clientUUID string) {
	h.Lock()
	defer h.Unlock()
	delete((*h.chats)[chatID], clientUUID)
	if len((*h.chats)[chatID]) == 0 {
		delete(*h.chats, chatID)
	}
}

func (h *HubScheduler) sendRoomMessages(msgType int, msgBytes []byte, chatID int, senderConn *websocket.Conn) {
	h.RLock()
	defer h.RUnlock()

	for _, conn := range (*h.chats)[chatID] {
		if conn == senderConn {
			continue
		}
		err := conn.WriteMessage(msgType, msgBytes)
		if err != nil {
			//h.removeClient(chatID, clientUUID)
			//cnErr := conn.Close()
			//if cnErr != nil {
			//	h.logger.Error(cnErr.Error())
			//}
			h.logger.Error(err.Error())
		}
	}
}

func (h *HubScheduler) listenRoomConnection(chatID int, clientUUID string, conn *websocket.Conn) {
	for {
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			h.removeClient(chatID, clientUUID)
			cnErr := conn.Close()
			if cnErr != nil {
				h.logger.Error(cnErr.Error())
			}
			h.logger.Error(err.Error())
			return
		}

		message := *models.InitMessageCreate(clientUUID, string(msgBytes), time.Now(), chatID)

		h.sendRoomMessages(msgType, msgBytes, chatID, conn)
		h.repoMessageCreator.CreateMessage(message)
	}
}

func (h *HubScheduler) JoinChat(chatID int, w http.ResponseWriter, r *http.Request) error {
	h.Lock()
	defer h.Unlock()

	if (*h.chats)[chatID] == nil {
		chatExists, err := h.repoChatGetter.IsChatExists(chatID)
		if err != nil {
			return err
		}
		if !chatExists {
			return RoomNotFound
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}

	clientUUID := uuid.New().String()

	switch (*h.chats)[chatID] {
	case nil:
		(*h.chats)[chatID] = map[string]*websocket.Conn{clientUUID: conn}
	default:
		(*h.chats)[chatID][clientUUID] = conn
	}

	go h.listenRoomConnection(chatID, clientUUID, conn)

	return nil
}
