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
		Short: "Run HTTP UI",
		Run: func(cmd *cobra.Command, args []string) {

			templatePath, err := cmd.Flags().GetString("template-path")
			if err != nil {
				log.Fatal(err)
			}

			addr, err := cmd.Flags().GetString("addr")
			if err != nil {
				log.Fatal(err)
			}

			s := server.New(db.Default, templatePath)
			err = s.Serve(addr)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().String("addr", ":53000", "http listen address")
	cmd.Flags().StringP("template-path", "t", "src/templates", "path to templates")
	return cmd
}
