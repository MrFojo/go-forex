package main

import (
	"github.com/mrfojo/go-forex/src/app"
)

func main() {
	app.EnsureInitializeData()
	app.Run()
}
