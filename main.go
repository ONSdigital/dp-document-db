package main

import (
	"fmt"
	"os"

	"github.com/ONSdigital/dp-document-db/database"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("application error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root := &cobra.Command{
		Use:   "db",
		Short: "TODO",
	}

	root.AddCommand(ping(), insert())

	return root.Execute()
}

func ping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the DocumentDB instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := getClient(cmd)
			if err != nil {
				return err
			}

			ctx, cancel := database.Ctx()
			defer cancel()

			err = cli.Ping(ctx, nil)
			if err != nil {
				return err
			}

			fmt.Println("connected to documentDB successful!")
			return nil
		},
	}

	cmd.Flags().StringP("username", "u", "", "DocumentDB username (required)")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringP("password", "p", "", "DocumentDB password (required)")
	cmd.MarkFlagRequired("password")

	return cmd
}

func insert() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insert",
		Short: "insert a value into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := getClient(cmd)
			if err != nil {
				return err
			}

			collection := cli.Database("metal-wiki").Collection("guitarists")

			ctx, cancel := database.Ctx()
			defer cancel()

			res, err := collection.InsertOne(ctx, bson.M{"name": "James Hetfield", "band": "Metallica"})
			if err != nil {
				return err
			}

			fmt.Printf("inserted 1 record ID: %s\n", res.InsertedID)
			return nil
		},
	}

	cmd.Flags().StringP("username", "u", "", "DocumentDB username (required)")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringP("password", "p", "", "DocumentDB password (required)")
	cmd.MarkFlagRequired("password")

	return cmd
}

func getClient(cmd *cobra.Command) (*mongo.Client, error) {
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return nil, err
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return nil, err
	}

	return database.NewClient(username, password)
}
