package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

type node struct {
	id string
	c  *websocket.Conn
}

type network struct {
	nodes   map[string]node
	addNode chan node
	rmNode  chan node
	bcChan  chan message
}

func (n *network) run() {
	for {
		select {
		case node := <-n.addNode:
			n.register(node)
		case node := <-n.rmNode:
			n.unregister(node)
		case mesg := <-n.bcChan:
			n.broadcastMesg(mesg)
		}
	}
}

func (n *network) register(node node) {
	n.nodes[node.c.RemoteAddr().String()] = node
}

func (n *network) unregister(node node) {
	delete(n.nodes, node.c.RemoteAddr().String())
}

func (n *network) broadcastMesg(m message) {
	for _, node := range n.nodes {
		if err := node.c.WriteJSON(m); err != nil {
			log.Println("Error broadcasting message:", err)
			return
		}
	}
}

func (n *network) spa(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", 302)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func (n *network) ws(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{}
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error accepting income connection:", err)
		return
	}

	node := node{id: r.URL.Query().Get("u"), c: wsc}
	n.addNode <- node

	go func() {
		for {
			var m message
			if err := node.c.ReadJSON(&m); err != nil {
				n.bcChan <- message{Text: fmt.Sprintf("%s disconnected", node.id), Sender: "Bot"}
				n.unregister(node)
				return
			}

			n.bcChan <- m
		}
	}()

}

func (n *network) connected(w http.ResponseWriter, r *http.Request) {
	var users []string

	for _, node := range n.nodes {
		users = append(users, node.id)
	}

	connectedNodes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(connectedNodes)
}

func Start(port string) error {
	n := &network{
		nodes:   make(map[string]node),
		addNode: make(chan node),
		rmNode:  make(chan node),
		bcChan:  make(chan message),
	}

	go n.run()

	mx := http.NewServeMux()
	fs := http.FileServer(http.Dir("node_modules"))

	mx.Handle("/node_modules/", http.StripPrefix("/node_modules/", fs))
	mx.HandleFunc("/", n.spa)
	mx.HandleFunc("/ws", n.ws)
	mx.HandleFunc("/chat/api/users", n.connected)

	return http.ListenAndServe(port, mx)
}
