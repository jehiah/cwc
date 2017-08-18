package main

import (
	"log"
	"strings"
	"time"

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
			base, err := cmd.Flags().GetString("base")
			if err != nil {
				log.Fatal(err)
			}
			if base == "" || !strings.HasSuffix(base, "/") {
				base = "/"
			}
			log.Printf("base URL is %s", base)

			readOnly, _ := cmd.Flags().GetBool("read-only")
			s := server.New(loadDB(cmd.Flags().GetString("db")), templatePath, base, readOnly)

			if skip, _ := cmd.Flags().GetBool("skip-browser-open"); !skip {
				go func() {
					time.Sleep(200 * time.Millisecond)
					err := s.OpenInBrowser()
					if err != nil {
						log.Println(err)
					}
				}()
			}

			err = s.Serve(addr)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	cmd.Flags().String("addr", ":53000", "http listen address")
	cmd.Flags().StringP("template-path", "t", "src/templates", "path to templates")
	cmd.Flags().String("base", "/", "Base URL Path")
	cmd.Flags().Bool("skip-browser-open", false, "skip oepening address in browser")
	cmd.Flags().Bool("read-only", false, "make UI read-only")
	return cmd
}
