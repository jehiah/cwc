package main

import (
	"log"

	"cwc/db"
	"cwc/server"

	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Web UI for viewing reports and adding notes",
		Run: func(cmd *cobra.Command, args []string) {

			templatePath, err := cmd.Flags().GetString("template-path")
			if err != nil {
				log.Fatal(err)
			}

			addr, err := cmd.Flags().GetString("addr")
			if err != nil {
				log.Fatal(err)
			}

			s := server.New(loadDB(cmd.Flags().GetString("db")), templatePath)
			err = s.Serve(addr)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	cmd.Flags().String("addr", ":53000", "http listen address")
	cmd.Flags().StringP("template-path", "t", "src/templates", "path to templates")
	return cmd
}
