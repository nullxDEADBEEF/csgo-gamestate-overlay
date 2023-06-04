package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	csgi "github.com/dank/go-csgsi"
)

var wss *websocket.Conn

type Message struct {
	Bomb  string
	Phase string
	Round string
	Map   string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// TODO: actually check the origin!
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade http connection to WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	wss = ws

	log.Println("Client Connected!")

	reader(ws)
}

// listens for new messages being sent to WebSocket endpoint
func reader(conn *websocket.Conn) {
	for {
		// read message
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("RECEIVED MSG FROM CLIENT: ", messageType, string(msg))

	}
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	game := csgi.New(10)

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		setupRoutes()
		http.ListenAndServe(":8080", nil)
		wg.Done()
	}()

	go func() {
		go func() {
			for state := range game.Channel {
				for weapon := range state.Player.Weapons {
					weapon := state.Player.Weapons[weapon]
					if weapon.State == "active" {
						fmt.Printf("%s clip: %d reserve: %d\n", weapon.Name, weapon.Ammo_clip, weapon.Ammo_reserve)
					}
				}
				fmt.Printf("Map: %s\n", state.Map.Name)
				fmt.Printf("WON ROUND: %s\n", state.Round.Win_team)
				fmt.Printf("ROUND BOMB: %s\n", state.Round.Bomb)
				fmt.Printf("MAP PHASE: %s\n", state.Map.Phase)

				if wss != nil {
					msg := Message{state.Round.Bomb, state.Round.Phase, state.Round.Win_team, state.Map.Name}
					msgJson, err := json.Marshal(msg)
					if err != nil {
						fmt.Println(err)
					}

					wss.WriteMessage(1, msgJson)
				}

			}
		}()

		fmt.Println("Server is ready for events!")

		game.Listen(":3000")
		wg.Done()
	}()

	wg.Wait()
}
