package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4"
	"github.com/netlify/netlify-commons/nconf"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Get the PostgreSQL URI from environment variable
	postgresqlURI := process.env.POSTGRESQL_URI
	if v, ok := request.QueryStringParameters["postgresqlURI"]; ok {
		postgresqlURI = v
	}

	// Create a connection pool using the PostgreSQL URI
	connConfig, err := pgx.ParseConfig(postgresqlURI)
	if err != nil {
		log.Printf("Failed to parse PostgreSQL URI: %v", err)
		return nil, err
	}
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: 5, // Adjust the number of connections as per your requirements
	})
	if err != nil {
		log.Printf("Failed to create connection pool: %v", err)
		return nil, err
	}
	defer connPool.Close()

	// Query the PostgreSQL version
	var version string
	err = connPool.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("Failed to query PostgreSQL version: %v", err)
		return nil, err
	}

	// Output the PostgreSQL version
	body := fmt.Sprintf("PostgreSQL version: %s", version)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       body,
	}, nil
}

func main() {
	// Initialize Netlify function configuration
	config, err := nconf.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Start the Lambda handler
	lambda.Start(handler)
}
