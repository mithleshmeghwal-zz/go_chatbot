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

	status := make(chan int)
	
	var callback models.Callback
	json.NewDecoder(r.Body).Decode(&callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				go func() {
					client := &http.Client{}
					log.Println("SENDER ID : ", event.Sender.ID)
					log.Println("Message Recieved: ", event.Message.Text)
					counter := 1
					value, err := redis.String(a.Conn.Do("GET", event.Sender.ID))
					if err == nil {
						// log.Println("GET", value)
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
					
					resp, err := client.Do(req)
					defer resp.Body.Close()
					if err != nil {
						log.Println("Message not Sent", err)
						status<- http.StatusInternalServerError
					} else {
						log.Println("Message Sent : ", strconv.Itoa(counter) + " ) " + event.Message.Text)
						ok, err := a.Conn.Do("SET", event.Sender.ID, counter)
						if ok != "OK" || err != nil {
							log.Println("Error: Failed to save value")
							status<- http.StatusInternalServerError
						} else {
							status<- http.StatusOK
						}
					}
				}()
			}
		}
		w.WriteHeader(<-status)
		w.Write([]byte("OK"))
		close(status)

	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}
}

func ProcessMessage(event models.Messaging, conn redis.Conn) {
	client := &http.Client{}
	log.Println("SENDER ID : ", event.Sender.ID)
	log.Println("Message Recieved: ", event.Message.Text)
	counter := 1
	value, err := redis.String(conn.Do("GET", event.Sender.ID))
	if err == nil {
		// log.Println("GET", value)
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
		log.Println("Error: Failed to save value")
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("Message not Sent", err)
	} else {
		log.Println("Message Sent : ", strconv.Itoa(counter) + " ) " + event.Message.Text)
	}
}
