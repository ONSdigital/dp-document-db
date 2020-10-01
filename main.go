package main

import (
	"context"
	"fmt"
	"os"

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
	return &cobra.Command{
		Use:   "docdb",
		Short: "TODO",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello world!")
			return nil
		},
	}
}
