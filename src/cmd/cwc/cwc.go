package main

import (
	"fmt"
	"log"
	"os"
	"testing/iotest"

	"cwc/db"
	"cwc/reporter"

	input "github.com/tcnksm/go-input"
)

func run(action string, args ...string) {
	var err error
	switch action {
	case "search":
		if len(args) >= 1 {
			search(args[0])
		} else {
			search("")
		}
	case "report":
		err = reporter.Run(db.Default, os.Stdout)
	case "new":
		err = newComplaint()
	case "short-reg", "short-regulations":
		listRegulations(true)
	case "reg", "regulations":
		listRegulations(false)
	case "help":
		fmt.Printf(`cwc -> Cyclists With Cameras

For more information see https://github.com/jehiah/cwc
`)
	default:
		log.Fatalf("not implemented")
	}
	if err != nil {
		log.Fatalf("%s", err)
	}
}

var ui *input.UI

func main() {
	ui = &input.UI{
		Reader: iotest.NewReadLogger("stdin:", os.Stdin),
	}

	if len(os.Args) > 1 {
		run(os.Args[1], os.Args[2:]...)
	} else {
		action, _ := ui.Select("", []string{"help", "search", "new", "report", "regulations", "short-regulations"}, &input.Options{
			Default: "new",
		})
		run(action)
	}
}

func search(query string) {
	var err error
	if query == "" {
		query, err = ui.Ask("Search for?", &input.Options{Required: true, Loop: true, HideOrder: true})
		if err != nil {
			os.Exit(1)
		}
		if query == "" {
			log.Fatalf("missing query")
		}
	} else {
		fmt.Printf("Searching for: %q\n", query)
	}

	files, err := db.Default.Find(query)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	for i, c := range files {
		if i == 0 {
			fmt.Printf("opening: %s %s\n", c, db.Default.FullPath(c))
			err = db.Default.ShowInFinder(c)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
			err = db.Default.Edit(c)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		} else {
			fmt.Printf("also found: %s %s\n", c, db.Default.FullPath(c))
		}
	}
}
