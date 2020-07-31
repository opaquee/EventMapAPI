package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
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

	observers := make(map[int](map[string]chan *model.Event), 1)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB:        db,
		Observers: observers,
	}}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: nil, //If CheckOrigin is nil: return false if the Origin request header is present and the origin host is not equal to request Host header.
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
