package main

import (
	"errors"
	"flag"
	"fmt"
	"goHtmlBuilder/builder"
	"goHtmlBuilder/fsPatrol"
	"goHtmlBuilder/httpServer"
	"goHtmlBuilder/minify"
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
	var isMinifyCss bool //TODO realize
	var isMinifyJs bool  //TODO realize

	flag.BoolVar(&isServe, "serve", false, "auto serve file updates, you don't need to reexecute ghtmler")
	flag.BoolVar(&isMinifyCss, "minifycss", false, "minify css, only for build, ignoring for serve")
	flag.BoolVar(&isMinifyJs, "minifyjs", false, "minify js, only for build, ignoring for serve")

	flag.Parse()

	necessaryFolders := []string{"dist", "ghtml", "liveReload", "static", "components"}
	for _, path := range necessaryFolders {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}

	State, errs := fsPatrol.GetState()
	if errs != nil {
		log.Fatal(errs)
	}

	err = builder.Build(State.GetGhtmlFiles(), isServe, minify.Params{
		IsMinifyCss: isMinifyCss,
		IsMinifyJs:  isMinifyJs,
	})
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
				err := builder.Build(currentState.GetGhtmlFiles(), isServe, minify.Params{
					IsMinifyCss: false,
					IsMinifyJs:  false,
				})
				if err != nil {
					log.Println(err)
				}
				ws_server.SendReload()
				prevState = currentState
			}

		}
	}

}
