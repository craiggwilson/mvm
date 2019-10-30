package main

import (
	"os"

	"github.com/craiggwilson/mvm/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
