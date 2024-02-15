package main

import (
	"fmt"
	"os"

	"github.com/ilya2049/gocomponent/internal/cli"
)

func main() {
	if err := cli.NewApp().Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
