package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	chatsAmount       = 20
	chatClientsAmount = 20 // will panic if less than 2
	serverUrl         = "95.84.137.217:3002"
	//serverUrl = "0.0.0.0:3002"
	messagesPerMinute    = 120 // in a minute :))))
	messagesSendDeadLine = 1   // minutes
)

type createChatResponse struct {
	ChatID int `json:"chat_id"`
}

// CreateChat returns id of created chat
func createChat() int {
	url := fmt.Sprintf("http://%v/chat/create_chat?name=%v", serverUrl, uuid.NewString())

	req, err := http.NewRequest("POST", url, nil) // No body is needed for this request
	if err != nil {
		log.Fatalf("Error creating request: %v", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err.Error())
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Status code != 200: %v, body: %v", resp.Status, string(body))
	}

	var chatResponse createChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		log.Fatalf("Error on umarshalling response: %v", err)
	}

	return chatResponse.ChatID
}

// 0.0.0.0:3002/ws/join_chat?id=2
func createChatClients(chatID int, output chan time.Time) {
	url := fmt.Sprintf("ws://%v/ws/join_chat?id=%v", serverUrl, chatID)

	var wg sync.WaitGroup

	for i := 0; i < chatClientsAmount; i++ {
		wg.Add(1)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2000)
			defer cancel()
			c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
			if err != nil {
				log.Fatalf("dial: %v", err.Error())
			}

			defer c.Close()

			wg.Done()

			for {
				mt, m, err := c.ReadMessage()
				if err != nil {
					log.Fatalf("read: %v", err.Error())
				}
				if mt == 1 || mt == 2 {
					output <- time.Now()
				} else {
					log.Fatalf("unexpected message type: %v, message: %v", mt, m)
				}
			}
		}()
	}

	wg.Wait()
}

type handledMessageData struct {
	chatID   int
	meanTime int
	maxTime  time.Duration
	minTime  time.Duration
}

func TestMessageLatency(t *testing.T) {
	chatsChannels := make(map[int]chan time.Time)
	for i := 0; i < chatsAmount; i++ {
		chatID := createChat()
		outputChan := make(chan time.Time, chatClientsAmount)
		chatsChannels[chatID] = outputChan
		createChatClients(chatID, outputChan)
	}

	startSyncChan := make(chan struct{})

	// магическое деление на 2, чтобы реально было нужное rpm, связано с тем, что отправитель ещё и обрабатывает ответы
	messageSendTicker := time.NewTicker(time.Minute / messagesPerMinute / 2)
	defer messageSendTicker.Stop()

	deadLineContext, cancel := context.WithTimeout(context.Background(), time.Minute*messagesSendDeadLine)
	defer cancel()

	var wg sync.WaitGroup

	for key := range chatsChannels {
		// if you're using go version older than 1.22!!!
		// key := key
		wg.Add(1)
		go func() {
			defer wg.Done()

			resultChan := make(chan handledMessageData, messagesPerMinute*messagesSendDeadLine)

			<-startSyncChan

			url := fmt.Sprintf("ws://%v/ws/join_chat?id=%v", serverUrl, key)
			c, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				log.Fatalf("dial: %v", err.Error())
			}

			defer c.Close()

		foreverLoop:
			for {
				select {
				case <-messageSendTicker.C:
					err = c.WriteMessage(1, []byte("hello, world!"))
					if err != nil {
						log.Fatalf("error on writing message to chatID: %v, err: %v ", key, err.Error())
					}

					timeStamp := time.Now()

					var meanDuration time.Duration
					var minTime = time.Hour
					var maxTime = time.Nanosecond

					var count = 0
					for elem := range chatsChannels[key] {
						duration := elem.Sub(timeStamp)
						meanDuration += duration
						minTime = min(minTime, duration)
						maxTime = max(maxTime, duration)
						count++
						if count == chatClientsAmount {
							break
						}
					}

					mean := int(meanDuration.Milliseconds()) / chatClientsAmount

					resultChan <- handledMessageData{
						chatID:   key,
						meanTime: mean,
						maxTime:  maxTime,
						minTime:  minTime,
					}
				case <-deadLineContext.Done():
					break foreverLoop
				}
			}

			for i := 0; i < len(resultChan); i++ {
				elem := <-resultChan
				fmt.Println(fmt.Sprintf("chatID: %v, mean time: %v milliseconds, max time: %v, min time: %v, "+
					"chatClients: %v", elem.chatID, elem.meanTime, elem.maxTime, elem.minTime, chatClientsAmount))
			}
		}()
	}

	close(startSyncChan)

	wg.Wait()

	fmt.Println("Done!")
}

//TODO: сваггер описания и сам сваггер для уже существующих хендлеров
//TODO: доделать хендлеры
//TODO: как замеряют пинг?
//TODO: в теле сообщения отправлять таймстамп отправки, тогда получатель знает время отправки и может просто записать всю эту инфу
//TODO: сделать профилирование в самом сервисе, а здесь только бомбер и затем считывание профайлинга с сервиса
