package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opaquee/EventMapAPI/graph"
	"github.com/opaquee/EventMapAPI/graph/generated"
	"github.com/opaquee/EventMapAPI/graph/model"
)

const defaultPort = "8080"

var db *gorm.DB

func main() {
	var err error

	log.Println("Connecting to database...")
	db, err := gorm.Open("postgres", "host=db port=5432 dbname=postgres user=user password=secret sslmode=disable")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	log.Println("Migrating tables...")
	db.AutoMigrate(&model.User{}, &model.Event{})

	log.Println("Starting server. Hold on to your potatoes!")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
