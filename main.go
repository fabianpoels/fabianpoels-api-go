package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabianpoels/fabianpoels-api-go/collections"
	"github.com/fabianpoels/fabianpoels-api-go/config"
	"github.com/fabianpoels/fabianpoels-api-go/db"
	"github.com/fabianpoels/fabianpoels-api-go/models"
	"github.com/fabianpoels/fabianpoels-api-go/server"
	"github.com/fabianpoels/fabianpoels-api-go/utils"
	"github.com/joho/godotenv"
)

func main() {
	environment := flag.String("e", "development", "environment")
	email := flag.String("email", "", "Valid email for user")
	passw := flag.String("password", "", "Password for user")
	role := flag.String("role", "user", "User role")
	os.Setenv("environment", *environment)
	task := flag.String("task", "server", "Task to run (server)")

	flag.Usage = func() {
		fmt.Println("Usage: go run main.go -e {mode} -task {server|add-user}")
	}
	flag.Parse()

	config.Init(*environment)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error reading env file. Err: %s", err)
	}

	switch *task {
	case "server":
		server.Init()
	case "db-indexes":
		db.CreateIndexes()
	case "add-user":
		addUser(email, passw, role)
	}

}

func addUser(email *string, passw *string, role *string) {
	if len(*email) < 3 {
		log.Fatalf("Email too short")
	}
	if len(*passw) < 3 {
		log.Fatalf("Passw too short")
	}

	passwHash, err := utils.HashPassword(*passw)
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		Email:     *email,
		Password:  passwHash,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err := collections.GetUserCollection(db.GetDbClient()).InsertOne(context.Background(), newUser)
	if err != nil {
		panic(err)
	}
	log.Println(res)
}
