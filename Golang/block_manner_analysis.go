package main

import (
	"flag"
	"fmt"
	"os"
)

type flags struct {
	mode *string
}

var cmdLineArgs flags

func init() {
	cmdLineArgs = flags{
		mode: flag.String("m", "select", "Select block mode in {chan, for, select}"),
	}
}

func fn_for() {
	for {
	}
}

func fn_select() {
	select {}
}

func fn_chan() {
	ch := make(chan struct{})
	<-ch
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	flag.Parse()

	fmt.Println("[main] begin...")

	go func() {
		fmt.Println("[goroutine] another goroutine...")
		for {
		}
	}()

	switch *cmdLineArgs.mode {
	case "select":
		fmt.Println("[main] choose select")
		fn_select()
	case "chan":
		fmt.Println("[main] choose chan")
		fn_chan()
	case "for":
		fmt.Println("[main] choose for")
		fn_for()
	default:
		fmt.Println("[main] can not be here")
	}
}
