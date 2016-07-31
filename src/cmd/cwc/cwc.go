package main

import (
	"fmt"
	"os"

	"db"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Printf("Searching for: %q\n", arg)
		files, err := db.Default.Find(arg)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
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
}
