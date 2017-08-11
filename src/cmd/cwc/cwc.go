package main

import (
	"fmt"
	"log"
	"os"

	"cwc/db"
	"cwc/reporter"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(listRegulations())
	RootCmd.AddCommand(newComplaint())
	RootCmd.AddCommand(serverCmd())
	RootCmd.AddCommand(report())
	RootCmd.AddCommand(json())
	RootCmd.AddCommand(editCmd())
	RootCmd.AddCommand(searchCmd())
}

var RootCmd = &cobra.Command{
	Use:   "cwc",
	Short: "Cyclists With Cameras",
	Long:  "Cyclists With Cameras - utilities for managing a database of T&LC complaints.\n\nFor more information see https://github.com/jehiah/cwc",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			os.Exit(1)
		}
	},
}

func report() *cobra.Command {
	cmd := &cobra.Command{
		Use: "report",
		Run: func(cmd *cobra.Command, args []string) {
			err := reporter.Run(db.Default, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}

func json() *cobra.Command {
	cmd := &cobra.Command{
		Use: "json",
		Run: func(cmd *cobra.Command, args []string) {
			body, err := reporter.JSON(db.Default)
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(body)
		},
	}
	return cmd
}

func editCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "edit",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				search(db.Default, args[0], "edit")
			} else {
				search(db.Default, "", "edit")
			}
		},
	}
	return cmd
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "search",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				search(db.Default, args[0], "search")
			} else {
				search(db.Default, "", "search")
			}
		},
	}
	return cmd
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
