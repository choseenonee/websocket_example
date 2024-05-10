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
	//serverUrl         = "url:3002"
	serverUrl = "0.0.0.0:3002"
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

func sendMessage(chatID int, message string) time.Time {
	url := fmt.Sprintf("ws://%v/ws/join_chat?id=%v", serverUrl, chatID)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("dial: %v", err.Error())
	}

	defer c.Close()

	err = c.WriteMessage(1, []byte(message))
	if err != nil {
		log.Fatalf("error on writing message to chatID: %v, err: %v ", chatID, err.Error())
	}

	return time.Now()
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

	var wg sync.WaitGroup

	for key := range chatsChannels {
		// if you're using go version older than 1.22!!!
		// key := key
		wg.Add(1)
		go func() {
			defer wg.Done()

			<-startSyncChan
			timeStamp := sendMessage(key, "hello, world!")

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
			fmt.Println(fmt.Sprintf("chatID: %v, mean time: %v milliseconds, max time: %v, min time: %v, chatClients: %v", key, mean, maxTime, minTime, chatClientsAmount))
		}()
	}

	close(startSyncChan)

	wg.Wait()

	fmt.Println("Done!")
}
