package configs

import (
	"flag"
	"goHtmlBuilder/minify"
)

type Config struct {
	IsServe bool
	Minify  minify.Config
}

func GetFlagsConfig() Config {
	var c Config

	var isOMCss bool

	flag.BoolVar(&c.IsServe, "serve", false, "auto serve file updates, you don't need to reexecute ghtmler")

	flag.BoolVar(&c.Minify.Css.IsMinify, "mcss", false, "minify css, only for build, ignoring for serve")
	flag.BoolVar(&c.Minify.Css.IsOptimize, "ocss", false, "optimize css, only for build, ignoring for serve")

	flag.BoolVar(&isOMCss, "omcss", false, "optimize and minify  css, only for build, ignoring for serve")

	flag.BoolVar(&c.Minify.Js.IsMinify, "mjs", false, "minify js, only for build, ignoring for serve")

	// TODO realize
	flag.BoolVar(&c.Minify.Js.IsOptimize, "ojs", false, "optimize js, only for build, ignoring for serve")

	flag.Parse()

	if isOMCss {
		c.Minify.Css.IsOptimize = true
		c.Minify.Css.IsMinify = true
	}

	return c
}
