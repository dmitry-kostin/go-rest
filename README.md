<!-- PROJECT LOGO -->
<br />
<p align="center">
  <img src="https://repository-images.githubusercontent.com/397879110/ca96c957-860d-4ec9-a37c-f3274b15d997" alt="Logo" width="640">

<h3 align="center">Go Rest API Starter Template</h3>
  <p align="center">
    A template for a fast start to building REST Web Services in Go. This project provides basic skeleton along with
    some of the best practices and tools for building RESTful APIs in Go.
    <br />
    <br />
    <a href="#some-of-the-features-included-in-this-template">Features</a>
    ·
    <a href="#getting-started">Getting Started</a>
    ·
    <a href="#authentication">Authentication</a>
    ·
    <a href="#error-handling">Error handling</a>
  </p>
</p>

## Brief introduction notes

Sometimes you just need to start a new project, but you don't want to spend time on setting up the project structure,
configuring the database, and do other things that are not strictly related to the business logic of your application
or idea you want to validate. 

This project aims to provide a good starting point for that purpose. It's not a framework, it's just a template where 
you can quickly remove the parts you don't need and add the ones you need.

### Some of the features included in this template:
- The project is configured to use basic manual dependency injection, which is a good fit for small to medium size projects
- PostgreSQL database with migrations, it's up to you on devops side how to run it on production
- Authentication with API keys, simple and yet effective
- Error handling, I really like idea of customizing error messages for clients
- Logging
- Configuration, while `.env` file is present in project root, it's up to you how to load it 
- Service layer, can serve as an excellent starting point for DDD services or basic CRUD operations
- Struct validations using govalidator as a fast and replaceble solution especially if you in favor of other methods


## Getting Started

### Prepare environment variables
* Create and fill the `.env` file according to the example `.example.env`
* Export environment variables from `.env` in any suitable way for you
  * eg.: `export $(grep -v '^#' .env | xargs)`
  * eg.: or set an alias for zshrc / bashrc `alias loadenv="export \$(grep -v '^#' .env | xargs)"`
  * eg.: utilities like [direnv](https://direnv.net/) also could be helpful

### Setup a database
* Install [docker](https://docs.docker.com/desktop/mac/install/)
* Run database image `docker compose up -d`
* Run migrations (refer to a section below)

### Migrations

Migrations are based on [golang-migrate](https://github.com/golang-migrate/migrate).

To install: `brew install golang-migrate`

Create a new migration:
```
migrate create -ext sql -dir db/migrations -seq create_users
```

Run migrations all the way up or down
```
migrate -database ${DB_URI} -path src/db/migrations up
migrate -database ${DB_URI} -path src/db/migrations down
```

### Run application in development mode

I use air for live-reloading during development

* Install [air](https://github.com/cosmtrek/air) for live reloading
* Run `air`

Or build and run the application directly `go build -o ./tmp ./cmd/api-service && ./tmp/api-service`

### Authentication

The application uses API keys for authentication, which are typically passed in the `Authorization` header using Bearer 
authentication scheme.

While this basic authentication may appear straightforward to implement, it's crucial to remain cognizant of potential
threats and attacks that could occur, particularly given the extensive privileges often associated with these keys.

Application middleware is focused on security preventing the most common attacks vectors: such as timing attacks and
key provisioning.

To generate API keys, you can use the following command:

`$ openssl rand -base64 32`

This command will generate a 32-byte random string, which is a good length for API keys
To be able to use that key in the application, you need to create a sha256 hash from it:

`$ echo -n "<you_key>" | sha256sum`

This command will generate a 32-byte hash, which you can put into the config `.env` file

`APP_API_KEYS` supports multiple keys, separated by comma, eg.: `key1,key2,key3`; This might be useful for keys 
rotation or when you want to provide granular access to different clients or revoke one.

### Error handling

The application uses [cockroachdb/errors](https://github.com/cockroachdb/errors) for error handling. This package 
provides enhanced error handling capabilities, such as wrapping errors, attaching stack traces, and adding context to
errors.

This allows us to provide more detailed human readable error messages to clients, while still maintaining the ability 
to log the original error internally.
