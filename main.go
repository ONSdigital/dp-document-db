package main

import (
	"fmt"
	"os"

	"github.com/ONSdigital/dp-document-db/v1"
	v2 "github.com/ONSdigital/dp-document-db/v2"
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

	root.AddCommand(v1Command(), v2Command())

	return root.Execute()
}

func v1Command() *cobra.Command {
	cmd := &cobra.Command{Use: "v1"}
	ping := addUsernamePasswordFlags(&cobra.Command{
		Use:   "ping",
		Short: "Ping the DocumentDB instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := getV1Client(cmd)
			if err != nil {
				return err
			}

			ctx, cancel := v1.Ctx()
			defer cancel()

			err = cli.Ping(ctx, nil)
			if err != nil {
				return err
			}

			fmt.Println("connected to documentDB successful!")
			return nil
		},
	})

	insert := addUsernamePasswordFlags(&cobra.Command{
		Use:   "insert",
		Short: "insert a value into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := getV1Client(cmd)
			if err != nil {
				return err
			}

			collection := cli.Database("metal-wiki").Collection("guitarists")

			ctx, cancel := v1.Ctx()
			defer cancel()

			res, err := collection.InsertOne(ctx, bson.M{"name": "James Hetfield", "band": "Metallica"})
			if err != nil {
				return err
			}

			fmt.Printf("inserted 1 record ID: %s\n", res.InsertedID)
			return nil
		},
	})

	cmd.AddCommand(ping, insert)
	return cmd
}

func addUsernamePasswordFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP("username", "u", "", "DocumentDB username (required)")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringP("password", "p", "", "DocumentDB password (required)")
	cmd.MarkFlagRequired("password")

	return cmd
}

func getV1Client(cmd *cobra.Command) (*mongo.Client, error) {
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return nil, err
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return nil, err
	}

	return v1.NewClient(username, password)
}

func v2Command() *cobra.Command {
	v2Cmd := &cobra.Command{
		Use:   "v2",
		Short: "use the global sign mongo lb to connect to the cluster",
	}

	ping := addUsernamePasswordFlags(&cobra.Command{
		Use:   "ping",
		Short: "ping the mongo cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, err := cmd.Flags().GetString("username")
			if err != nil {
				return err
			}

			password, err := cmd.Flags().GetString("password")
			if err != nil {
				return err
			}

			session, err := v2.NewSession(username, password)
			if err != nil {
				return err
			}

			defer session.Close()

			err = session.Ping()
			if err != nil {
				return err
			}

			fmt.Println("mgo global sign ping successful")

			dbs, err := session.DatabaseNames()
			if err != nil {
				return err
			}

			fmt.Printf("dbs: %+v\n", dbs)
			return nil
		},
	})

	v2Cmd.AddCommand(ping)
	return v2Cmd
}
