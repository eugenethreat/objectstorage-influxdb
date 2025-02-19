package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// https://pkg.go.dev/github.com/influxdata/influxdb-client-go/v2#section-readme

func main() {
	// load env. variables
	err := godotenv.Load(".env")

	fmt.Println("doing it")
	INFLUX_URL := os.Getenv("INFLUX_URL")
	INFLUX_TOKEN := os.Getenv("INFLUX_TOKEN")

	client := influxdb2.NewClient(INFLUX_URL, INFLUX_TOKEN)

	// validate client connection health
	health, err := client.Health(context.Background())
	if err != nil {
		log.Fatalf("error connecting to db")
	} else {
		fmt.Println("successfully connected to db - STATUS: " + health.Status)
	}

}
