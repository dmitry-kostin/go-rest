<!-- PROJECT LOGO -->
<br />
<p align="center">
  <img src="https://www.machinelearningmastery.ru/img/0-163193-239529.png" alt="Logo" width="450" height="173">

<h3 align="center">Go Rest API Starter Template</h3>
  <p align="center">
    A template for a fast start to building REST Web Services in Go. Uses gorilla/mux as a router/dispatcher and Negroni as a middleware handler.
    <br />
    <br />
    <a href="#getting-started">Getting Started</a>
    ·
    <a href="#migrations">Migrations</a>
    ·
    <a href="#contributing">Contributing</a>
  </p>
</p>


## Getting Started

### Prepare environment variables
* Create and fill the `.env` file according to the example `.example.env`
* Export environment variables in a single command `export $(grep -v '^#' .env | xargs)`

### Setup a database
* Install [docker](https://docs.docker.com/desktop/mac/install/)
* Run database image `docker-compose up -d --build`
* Run migrations (refer to section below)

### Run application in development mode

We use air for live-reloading during development

* Install [air](https://github.com/cosmtrek/air) for live reloading 
* Run `air`

### Build and run

* `go build -o ./tmp ./cmd/api-service`
* `./tmp/api-service`

### Migrations

Migrations are based on [golang-migrate](https://github.com/golang-migrate/migrate).

To install: `brew install golang-migrate`

Create a new migration:
```
migrate create -ext sql -dir db/migrations -seq create_users
```

Run migrations all the way up or down
```
migrate -database ${DB_URI} -path ./db/migrations up
migrate -database ${DB_URI} -path db/migrations down
```