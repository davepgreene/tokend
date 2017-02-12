package main

import (
	"fmt"
	"os"

	"github.com/davepgreene/tokend/cmd"
)

func main() {
	if err := cmd.TokendCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
