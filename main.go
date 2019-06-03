package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"x_yield/cmd/cli"
)

func main() {
	flag.Usage = func() { cli.UsageExit(0) }
	flag.Parse()
	args := flag.Args()

}
