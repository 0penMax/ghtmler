package main

import (
	"errors"
	"fmt"
	"goHtmlBuilder/builder"
	"goHtmlBuilder/configs"
	"goHtmlBuilder/fsPatrol"
	"goHtmlBuilder/httpServer"
	ws_server "goHtmlBuilder/wsServer"
	"log"
	"os"
	"time"
)

func main() {
	//TODO think about system where css and js incude inside html file
	//TODO think about system where all images from page optimize by size and transform in sprite(one big image)
	f, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	log.SetOutput(f)

	createNecessaryFolders()

	State, errs := fsPatrol.GetState()
	if errs != nil {
		log.Fatal(errs)
	}

	//TODO create additional json config
	config := configs.GetFlagsConfig()

	err = builder.Build(State.GetGhtmlFiles(), config.IsServe, config.Minify)
	if err != nil {
		fmt.Println(err)
		return
	}

	if config.IsServe {
		fmt.Println("serve...")
		ws_server.StartServer()
		httpServer.RunServer()
		prevState := State
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			time.Sleep(1 * time.Second)
			currentState, err := fsPatrol.GetState()
			if err != nil {
				fmt.Println(err)
				return
			}

			if fsPatrol.IsDiffState(prevState, currentState) {
				fmt.Println("rebuild")
				err := builder.Build(currentState.GetGhtmlFiles(), config.IsServe, config.Minify)
				if err != nil {
					fmt.Println(err)
					return
				}
				ws_server.SendReload()
				prevState = currentState
			}

		}
	}

}

func createNecessaryFolders() {
	fmt.Println("check necessary folder")
	necessaryFolders := []string{"dist", "ghtml", "liveReload", "static", "components"}
	for _, path := range necessaryFolders {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("creating %s folder\n", path)
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
