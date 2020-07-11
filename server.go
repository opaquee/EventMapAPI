package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opaquee/EventMapAPI/graph"
	"github.com/opaquee/EventMapAPI/graph/generated"
	"github.com/opaquee/EventMapAPI/graph/model"
	"github.com/opaquee/EventMapAPI/helpers/auth"
	"github.com/opaquee/EventMapAPI/helpers/dbconn"
)

var db *gorm.DB

const defaultPort = "8080"

func main() {
	time.Sleep(10 * time.Second)

	log.Println("Connecting to database...")
	db, err := dbconn.Open()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	log.Println("Migrating tables...")
	if err := db.AutoMigrate(&model.User{}, &model.Event{}).Error; err != nil {
		panic(err)
	}
	if err := db.Model(&model.Event{}).AddForeignKey("owner_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		panic(err)
	}

	log.Println("Starting server. Hold on to your potatoes!")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("Applying middleware...")
	router := chi.NewRouter()
	router.Use(auth.Middleware(db))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
