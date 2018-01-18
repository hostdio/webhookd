package main

import (
	"log"
	"os"

	"github.com/hostdio/webhookd/httpapi"
	_ "github.com/lib/pq"
)

var (
	postgresConnString = os.Getenv("POSTGRES_CONNECTION_STRING")
)

func main() {

	if postgresConnString == "" {
		panic("POSTGRES_CONNECTION_STRING is missing")
	}

	log.Fatal(httpapi.Cmd(postgresConnString))

}
