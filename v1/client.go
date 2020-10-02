package v1

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ONSdigital/dp-document-db/certs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Path to the AWS CA file
	caFilePath = "rds-combined-ca-bundle.pem"

	// Timeout operations after N seconds
	connectTimeout = 5

	// Timeout queries after N seconds
	queryTimeout = 5

	// Which instances to read from
	readPreference = "secondaryPreferred"

	connectionStringTemplate = "mongodb://%s:%s@%s/?ssl=true&ssl_ca_certs=rds-combined-ca-bundle.pem&replicaSet=rs0&readPreference=%s&retryWrites=false"
)

type ClientErr struct {
	Cause   error
	Message string
}

func (e ClientErr) Error() string {
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
}

func NewClient(username, password string) (*mongo.Client, error) {
	clusterEndpoint := os.Getenv("DB_ENDPOINT")

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint, readPreference)

	tlsConfig, err := certs.GetCustomTLSConfig(caFilePath)
	if err != nil {
		return nil, ClientErr{Message: "Failed getting TLS configuration", Cause: err}
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	if err != nil {
		return nil, ClientErr{Message: "Failed to create client", Cause: err}
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, ClientErr{Message: "Failed to connect to cluster", Cause: err}
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, ClientErr{Message: "Failed to ping cluster: %v", Cause: err}
	}

	fmt.Println("successfully connected to DocumentDB")
	return client, nil
}

func Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*queryTimeout)
}