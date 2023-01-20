package main

import (
	"fmt"
	"log"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/input"
)

func search(d db.ReadWrite, query, action string) {
	var err error
	if query == "" {
		query, err = input.Ask("Search for?", "")
		if err != nil {
			log.Fatalf("%s", err)
		}
		if query == "" {
			log.Fatalf("missing query")
		}
	} else {
		fmt.Printf("Searching for: %q\n", query)
	}

	files, err := d.Find(query)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	if id, ok := d.(db.Interactive); ok {
		for i, c := range files {
			if i == 0 {
				fmt.Printf("opening: %s %s\n", c, d.FullPath(c))
				switch action {
				case "search":
					err = id.ShowInFinder(c)
					if err != nil {
						fmt.Printf("%s\n", err)
					}
				}
				err = id.ShowInEditor(c)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			} else {
				fmt.Printf("also found: %s %s\n", c, d.FullPath(c))
			}
		}
	}
}
