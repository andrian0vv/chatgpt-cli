package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrian0vv/chatgpt-cli/cmd"
)

func main() {
	ctx := context.Background()

	if err := cmd.Execute(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
