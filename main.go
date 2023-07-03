package main

import (
	"flag"
	"fmt"
	"goHtmlBuilder/builder"
	"goHtmlBuilder/fsPatrol"
	"log"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var isServe bool

	flag.BoolVar(&isServe, "serve", false, "auto serve file updates, you don't need to reexecute ghtmler")

	flag.Parse()

	State, errs := fsPatrol.GetState()
	if errs != nil {
		log.Fatal(err)
	}

	err = builder.Build(State.GetGhtmlFiles())
	if err != nil {
		log.Fatal(err)
	}

	if isServe {
		fmt.Println("serve...")
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
				err := builder.Build(currentState.GetGhtmlFiles())
				if err != nil {
					log.Println(err)
				}
				prevState = currentState
			}

		}
	}

}
