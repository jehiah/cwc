package main

import (
	"fmt"
	"log"

	"cwc/reg"
	"github.com/spf13/cobra"
)

func listRegulations() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reg -s",
		Aliases: []string{"regulations", "list-regulations"},
		Short:   "List Regulations",
		Run: func(cmd *cobra.Command, args []string) {
			isShort, err := cmd.Flags().GetBool("short")
			if err != nil {
				log.Fatal(err)
			}
			for _, r := range reg.All {
				desc := r.Description
				if isShort && r.Short != "" {
					desc = r.Short
				}
				fmt.Printf("%s,%s\n", r.Code, desc)
			}
		},
	}
	cmd.Flags().BoolP("short", "s", false, "short format output")
	return cmd
}
