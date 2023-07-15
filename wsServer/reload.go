//base on https://github.com/markvincze/golang-reload-browser

package ws_server

import (
	"bytes"
	"log"
	"net/http"
)

var (
	hub *Hub
	// The port on which we are hosting the reload server has to be hardcoded on the client-side too.
	reloadAddress = ":12450"
)

func StartServer() {
	hub = newHub()
	go hub.run()

	mux := http.NewServeMux()

	mux.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go startServer(mux)
	log.Println("Reload server listening at", reloadAddress)
}

func startServer(mux *http.ServeMux) {
	server := http.Server{
		Addr:    reloadAddress,
		Handler: mux,
	}
	err := server.ListenAndServe()

	if err != nil {
		log.Println("Failed to start up the Reload server: ", err)
		return
	}
}

func SendReload() {
	message := bytes.TrimSpace([]byte("reload"))
	hub.broadcast <- message
}
