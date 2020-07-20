# EventMapAPI
Web API for the EventMap website. Users can log in, create, read, update, and delete event information. A client can also subscribe to realtime updates of new event data.

# Installation
Install these:
Docker: 
https://docs.docker.com/get-docker/
Golang:
https://golang.org/doc/install

# Running the project
first, clone this repo to your machine and open a terminal in the project directory

$ go generate ./...

$ docker-compose up --build


if you are on linux, instead run the second command like this:
$ sudo docker-compose up --build

# Sending Requests and Viewing Database
go to localhost:8080 in your browser to send requests to the API.



