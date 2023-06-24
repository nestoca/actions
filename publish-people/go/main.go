package main

import (
	"fmt"
	"os"

	"github.com/nestoca/actions/publish-people/go/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
