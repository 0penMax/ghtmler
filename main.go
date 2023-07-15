package main

import (
	"flag"
	"fmt"
	"goHtmlBuilder/builder"
	"goHtmlBuilder/fsPatrol"
	"goHtmlBuilder/httpServer"
	ws_server "goHtmlBuilder/wsServer"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("start")
	f, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var isServe bool

	flag.BoolVar(&isServe, "serve", false, "auto serve file updates, you don't need to reexecute ghtmler")

	flag.Parse()

	fmt.Println("flag parsed")

	State, errs := fsPatrol.GetState()
	if errs != nil {
		log.Fatal(errs)
	}

	err = builder.Build(State.GetGhtmlFiles(), isServe)
	if err != nil {
		log.Fatal(err)
	}

	if isServe {
		fmt.Println("serve...")
		ws_server.StartServer()
		httpServer.RunServer()
		prevState := State
		if err != nil {
			log.Fatal(err)
		}
		for {
			time.Sleep(1 * time.Second)
			currentState, err := fsPatrol.GetState()
			if err != nil {
				log.Println(err)
			}

			if fsPatrol.IsDiffState(prevState, currentState) {
				fmt.Println("rebuild")
				err := builder.Build(currentState.GetGhtmlFiles(), isServe)
				if err != nil {
					log.Println(err)
				}
				ws_server.SendReload()
				prevState = currentState
			}

		}
	}

}
