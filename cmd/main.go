package main

import (
	"qimenpan/internal/router"
)

func main() {
	engine := router.InitRouter()
	err := engine.Run(":8888")
	if err != nil {
		panic(err)
	}
}
