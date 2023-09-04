package main

import (
	"github.com/Smylet/symlet-backend/utilities/cli/root"
)

func main() {
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
