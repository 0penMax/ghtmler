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

	var isServe bool
	var isMinifyCss bool
	var isOptimizeCss bool
	var isOMCss bool

	var isOptimizeJs bool
	var isMinifyJs bool

	flag.BoolVar(&isServe, "serve", false, "auto serve file updates, you don't need to reexecute ghtmler")

	flag.BoolVar(&isMinifyCss, "minifycss", false, "minify css, only for build, ignoring for serve")
	flag.BoolVar(&isMinifyCss, "mcss", false, "minify css, only for build, ignoring for serve")

	flag.BoolVar(&isOptimizeCss, "optimizecss", false, "optimize css, only for build, ignoring for serve")
	flag.BoolVar(&isOptimizeCss, "ocss", false, "optimize css, only for build, ignoring for serve")

	flag.BoolVar(&isOMCss, "omcss", false, "optimize and minify  css, only for build, ignoring for serve")

	flag.BoolVar(&isMinifyJs, "minifyjs", false, "minify js, only for build, ignoring for serve")
	flag.BoolVar(&isMinifyJs, "mjs", false, "minify js, only for build, ignoring for serve")

	// TODO realize
	flag.BoolVar(&isOptimizeJs, "optimizejs", false, "optimize js, only for build, ignoring for serve")
	flag.BoolVar(&isOptimizeJs, "ojs", false, "optimize js, only for build, ignoring for serve")

	flag.Parse()

	createNecessaryFolders()

	State, errs := fsPatrol.GetState()
	if errs != nil {
		log.Fatal(errs)
	}

	var p minify.Params

	if isOMCss {
		p = minify.Params{
			IsMinifyCss:   true,
			IsMinifyJs:    isMinifyJs,
			IsOptimizeCss: true,
			IsOptimizeJs:  isOptimizeJs,
		}
	} else {
		p = minify.Params{
			IsMinifyCss:   isMinifyCss,
			IsMinifyJs:    isMinifyJs,
			IsOptimizeCss: isOptimizeCss,
			IsOptimizeJs:  isOptimizeJs,
		}
	}

	err = builder.Build(State.GetGhtmlFiles(), isServe, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isServe {
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
				err := builder.Build(currentState.GetGhtmlFiles(), isServe, minify.Params{
					IsMinifyCss: false,
					IsMinifyJs:  false,
				})
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
