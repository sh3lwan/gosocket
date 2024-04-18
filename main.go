package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clients = make(map[string]Client)

var broadcast = make(chan *ReceivedMessage)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// Handle client messages
	go func() {
		for {
			sender := <-broadcast

			fmt.Printf("client count: %d\n", len(clients))

			for _, client := range clients {
				conn := client.Conn

				message, err := json.Marshal(ReceivedMessage{
					Id:      sender.Id,
					Message: sender.Message,
					IsNew:   sender.IsNew,
				})

				getMessagesCollection().InsertOne(context.Background(), message)

				if err != nil {
					log.Printf("Error writing message: %v", err)
				}

				err = conn.WriteMessage(websocket.TextMessage, message)

				fmt.Printf("message sent: %s %s\n", client.Id, message)

				if err != nil {
					log.Printf("Error writing message: %v", err)
				}

			}
		}

	}()

	http.HandleFunc("/", wsHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Works"))
	})

	http.HandleFunc("/api/connect", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("something went wrong: %v", err.Error())
		}
		var authBody AuthBody
		err = json.Unmarshal(body, &authBody)
		if err != nil {
			fmt.Printf("something went wrong: %v", err.Error())
		}

		messages := getMessages()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			StatusCode: http.StatusOK,
			Data: map[string]any{
				"messages": messages,
			},
		})

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Response struct {
	StatusCode int16 `json:"status_code"`
	Data       map[string]any
}
type AuthBody struct {
	Username string `json:"username"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

type Client struct {
	Id string
	*websocket.Conn
}

type ReceivedMessage struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	Username string `json:"username"`
	IsNew    bool   `json:"is_new"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("something went wrong %v\n", err)
		return
	}

	uuid := uuid.New().String()

	client := &Client{
		Conn: conn,
		Id:   uuid,
	}

	clients[uuid] = *client

	fmt.Printf("new connection established: %s\n", uuid)

	fmt.Printf("clients count: %d\n", len(clients))

	go handleMessages(client)
}

func handleMessages(client *Client) {
	broadcast <- &ReceivedMessage{
		Id:      client.Id,
		Message: "Connection established",
		IsNew:   true,
	}

	for {
		_, message, err := client.Conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:\n", err)
			break
		}

		broadcast <- &ReceivedMessage{
			Id:      client.Id,
			Message: string(message),
		}
	}

	fmt.Printf("Connection Closed: %s\n", client.Id)

	client.Conn.Close()

	delete(clients, client.Id)
}

func getDB() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://root:password@localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return client
}

func getMessagesCollection() *mongo.Collection {
	client := getDB()

	collection := client.Database("chat").Collection("messages")

	return collection
}

func getMessages() []ReceivedMessage {
	filter := bson.M{}

	options := options.Find()

	cursor, err := getMessagesCollection().Find(context.Background(), filter, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var messages []ReceivedMessage

	// Iterate over the cursor and process each document
	for cursor.Next(context.Background()) {
		var message ReceivedMessage
		if err := cursor.Decode(&message); err != nil {
			log.Fatal(err)
		}
		// Process the message document
		log.Println(message)
		messages = append(messages, message)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return messages
}
