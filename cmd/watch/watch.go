// watch is a tiny command-line driver that exercises github.com/cespare/fswatch.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/cespare/fswatch"
)

func main() {
	coalesceInterval := flag.Duration("d", 500*time.Millisecond, "Coalesce interval")
	flag.Parse()

	var dir string
	switch flag.NArg() {
	case 0:
		dir = "."
	case 1:
		dir = flag.Arg(0)
	default:
		log.Fatal("usage: watch [flags] [dir]")
	}

	events, errors, err := fswatch.Watch(dir, *coalesceInterval)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Fatal(<-errors)
	}()
	for ev := range events {
		fmt.Println(ev)
	}
}
