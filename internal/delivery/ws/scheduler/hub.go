package scheduler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
	"websockets/internal/metrics"
	"websockets/internal/models"
	"websockets/internal/repository"
	"websockets/pkg/log"
)

func InitHubScheduler(logger *log.Logs, repoMessageCreator RepoMessageCreator,
	repoChatGetter repository.ChatGetterRepo, prometheusMetrics metrics.PrometheusMetrics) *HubScheduler {
	return &HubScheduler{
		logger:             logger,
		chats:              &map[int]map[string]*websocket.Conn{},
		repoMessageCreator: repoMessageCreator,
		repoChatGetter:     repoChatGetter,
		prometheusMetrics:  prometheusMetrics,
	}
}

type HubScheduler struct {
	sync.RWMutex
	logger             *log.Logs
	chats              *map[int]map[string]*websocket.Conn
	repoChatGetter     repository.ChatGetterRepo
	repoMessageCreator RepoMessageCreator
	prometheusMetrics  metrics.PrometheusMetrics
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *HubScheduler) removeClient(chatID int, clientUUID string) {
	h.Lock()
	defer h.Unlock()

	delete((*h.chats)[chatID], clientUUID)
	h.prometheusMetrics.UsersOnline.Dec()

	if len((*h.chats)[chatID]) == 0 {
		delete(*h.chats, chatID)
		h.prometheusMetrics.ChatsOnline.Dec()
	}
}

func (h *HubScheduler) sendRoomMessages(msgType int, message *models.MessageCreate, senderConn *websocket.Conn) {
	h.RLock()
	defer h.RUnlock()

	for _, conn := range (*h.chats)[message.ChatID] {
		if conn == senderConn {
			continue
		}
		jsonBytes, _ := json.Marshal(message)
		err := conn.WriteMessage(msgType, jsonBytes)
		if err != nil {
			h.logger.Error(err.Error())
		}
	}
	h.prometheusMetrics.MessagesLatency.Observe(time.Since(message.SendTimeStamp).Seconds())
}

func (h *HubScheduler) listenRoomConnection(chatID int, clientUUID string, conn *websocket.Conn) {
	defer conn.Close()
	for {
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			h.removeClient(chatID, clientUUID)
			h.logger.Error(err.Error())
			return
		}

		message := models.InitMessageCreate(clientUUID, string(msgBytes), time.Now(), chatID)

		h.sendRoomMessages(msgType, message, conn)
		h.prometheusMetrics.MessagesSent.Inc()
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

	h.prometheusMetrics.UsersOnline.Inc()

	switch (*h.chats)[chatID] {
	case nil:
		(*h.chats)[chatID] = map[string]*websocket.Conn{clientUUID: conn}
		h.prometheusMetrics.ChatsOnline.Inc()
	default:
		(*h.chats)[chatID][clientUUID] = conn
	}

	go h.listenRoomConnection(chatID, clientUUID, conn)

	return nil
}
