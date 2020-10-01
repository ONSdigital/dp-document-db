package main

import (
	"context"
	"os"
	"time"

	"github.com/ONSdigital/dp-document-db/database"
	"github.com/ONSdigital/log.go/log"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		log.Event(context.Background(), "application error", log.ERROR, log.Error(err))
		os.Exit(1)
	}
}

func run() error {
	root := &cobra.Command{
		Use:   "docdb",
		Short: "TODO",
	}

	root.AddCommand(ping())

	return root.Execute()
}

func ping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the DocumentDB instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, err := cmd.Flags().GetString("username")
			if err != nil {
				return err
			}

			password, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			}

			cli, err := database.NewClient(username, password)
			if err != nil {
				return err
			}

			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			defer cli.Disconnect(ctx)
			return nil
		},
	}

	cmd.Flags().StringP("username", "u", "", "DocumentDB username (required)")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringP("password", "p", "", "DocumentDB password (required)")
	cmd.MarkFlagRequired("password")

	return cmd
}
