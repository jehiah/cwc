package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/internal/reporter"
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

func loadDB(p string, err error) db.ReadWrite {
	if p == "" {
		return db.Default
	}
	return db.LocalFilesystem(p)
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
		Use:   "report",
		Short: "Text format summarized view of report activity",
		Run: func(cmd *cobra.Command, args []string) {
			err := reporter.Run(loadDB(cmd.Flags().GetString("db")), os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	return cmd
}

func json() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "Output all complaints as JSON",
		Run: func(cmd *cobra.Command, args []string) {
			db := loadDB(cmd.Flags().GetString("db"))
			err := reporter.JSON(os.Stdout, db)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	return cmd
}

func editCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit complaint notes.txt",
		Example: "edit [query]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				search(loadDB(cmd.Flags().GetString("db")), args[0], "edit")
			} else {
				search(loadDB(cmd.Flags().GetString("db")), "", "edit")
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	return cmd
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search",
		Short:   "Search for a complaint by keword",
		Example: "search [query]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) >= 1 {
				search(loadDB(cmd.Flags().GetString("db")), args[0], "search")
			} else {
				search(loadDB(cmd.Flags().GetString("db")), "", "search")
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	return cmd
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
