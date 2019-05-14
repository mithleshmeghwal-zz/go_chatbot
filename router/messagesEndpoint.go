package router

import (
	"net/http"
	"log"
	"chatbot/models"
	"encoding/json"
	"bytes"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"os"
	"fmt"
)

const (
	FACEBOOK_API = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

func (a *Router)MessagesEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Message Event")
	var callback models.Callback
	json.NewDecoder(r.Body).Decode(&callback)
	log.Println("MESSAGE RECIEVED ", callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				ProcessMessage(event, a.Conn)
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}
}

func ProcessMessage(event models.Messaging, conn redis.Conn) {
	client := &http.Client{}
	log.Println("SENDER", event.Sender.ID, event.Message.Text)
	counter := 1
	value, err := redis.String(conn.Do("GET", event.Sender.ID))
	if err == nil {
		log.Println("GET", value)
		counter, err = strconv.Atoi(value)
		if err != nil {
			log.Println("Corrupted Data")
		} else {
			counter++;
		}
	}

	response := models.Response{
		Recipient: models.User{
			ID: event.Sender.ID,
		},
		Message: models.Message{
			Text: strconv.Itoa(counter) + " ) " + event.Message.Text,
		},
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)
	url := fmt.Sprintf(FACEBOOK_API, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	ok, err := conn.Do("SET", event.Sender.ID, counter)
	if ok != "OK" || err != nil {
		log.Println("Failed to sa ve value")
	}
	log.Println("SET", ok, err)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	defer resp.Body.Close()
}
