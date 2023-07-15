package httpServer

import (
	"goHtmlBuilder/builder"
	"log"
	"net/http"
)

func RunServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("dist/static"))))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(builder.LIVE_RELOAD_FOLDER))))

	log.Println("Listening...")

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

}
