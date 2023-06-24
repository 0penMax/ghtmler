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
	if isServe {
		fmt.Println("serve...")
		prevState, err := fsPatrol.GetState()
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
				err := builder.Build()
				if err != nil {
					log.Println(err)
				}
				prevState = currentState
			}

		}
	} else {
		err := builder.Build()
		if err != nil {
			log.Fatal(err)
		}
	}

}
