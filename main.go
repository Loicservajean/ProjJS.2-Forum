package main

import (
	"rompelago/app"
)

func main() {
	application := app.InitApp()
	defer application.Close()
	application.Start()
}
