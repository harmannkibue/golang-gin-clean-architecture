(docs/img/golang-gopher.png)

# Golang Clean template

Clean Architecture template for Golang Microservices

[![License](https://img.shields.io/github/license/evrone/go-clean-template.svg)](https://github.com/evrone/go-clean-template/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/evrone/go-clean-template.svg)](https://github.com/harmannkibue/golang_gin_clean_architecture/releases/)
[![codecov](https://codecov.io/gh/evrone/go-clean-template/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/evrone/go-clean-template)

## Overview
This project shows how to organise golang microservice code according to Robert Martin (aka Uncle Bob).
It aims to :
- Show how to organise the codebase to avoid supergheti.
- Show the code generation tooling for efficiency. 
- Demonstration of 
- Demonstrate a scalable and easy to maintain golang structure

## Technology/Tooling
- [Gin gonic](https://github.com/gin-gonic/gin) http framework
- [SQLC](https://docs.sqlc.dev/en/latest/index.html) Compile SQL to type-safe code
- [Mockery](https://github.com/vektra/mockery) mock testing framework
- [Swag](https://github.com/swaggo/swag) generating swagger documentation
- [Golang migrate](https://github.com/golang-migrate/migrate) for running database migrations

## Content
- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [Dependency Injection](#dependency-injection)
- [Clean Architecture](#clean-architecture)

## Quick start
Seeing all Configuration commands:
```azure
$ make 
```
Local development:
```sh
# Run app with postgres integration
$ make composeUp
```

Testing with Mockery:
```sh
# Run test with mocks for database
$ make testWithCoverProfile
```

## Project structure
### `cmd/app/main.go`
Configuration and logger initialization. Then the main function "continues" in
`internal/app/app.go`.

### `config`
Configuration. First, `config.yml` is read, then environment variables overwrite the yaml config if they match.
The config structure is in the `config.go`.
The `env-required: true` tag obliges you to specify a value (either in yaml, or in environment variables).

For configuration, we chose the [cleanenv](https://github.com/ilyakaznacheev/cleanenv) library.
It does not have many stars on GitHub, but is simple and meets all the requirements.

Reading the config from yaml contradicts the ideology of 12 factors, but in practice, it is more convenient than
reading the entire config from ENV.
It is assumed that default values are in yaml, and security-sensitive variables are defined in ENV.

### `docs`
Swagger documentation. Auto-generated by [swag](https://github.com/swaggo/swag) library.
You don't need to correct anything by yourself.

### `internal/app`
There is always one _Run_ function in the `app.go` file, which "continues" the _main_ function.

This is where all the main objects are created.
Dependency injection occurs through the "New ..." constructors (see Dependency Injection).
This technique allows us to layer the application using the [Dependency Injection](#dependency-injection) principle.
This makes the business logic independent from other layers.

Next, we start the server and wait for signals in _select_ for graceful completion.
If `app.go` starts to grow, you can split it into multiple files.

The `migrate.go` file is used for database auto migrations.
It is included if an argument with the _migrate_ tag is specified.
For example:

```sh
$ go run -tags migrate ./cmd/app
```

### `internal/controller`
Server handler layer (MVC controllers). Only http server implemented GRPC later:
- REST http (Gin framework)

Server routers are written in the same style:
- Handlers are grouped by area of application (by a common basis)
- For each group, its own router structure is created, the methods of which process paths
- The structure of the business logic is injected into the router structure, which will be called by the handlers

#### `internal/controller/http`
Simple REST versioning.
For v2, we will need to add the `http/v2` folder with the same content.
And in the file `internal/app` add the line:

```go
handler := gin.New()
v1.NewRouter(handler, t)
v2.NewRouter(handler, t)
```

Instead of Gin, you can use any other http framework or even the standard `net/http` library.

In `v1/router.go` and above the handler methods, there are comments for generating swagger documentation using [swag](https://github.com/swaggo/swag).

### `internal/entity`
This contains items that are accessible from any file. e.g Interfaces, test mocks etc

### `internal/usecase`
Business logic.
- Methods are grouped by area of application (on a common basis)
- Each group has its own structure
- One file - one structure

#### `internal/usecase/repositories`
A repository is an abstract storage (database) that business logic works with. We use golang [SQLC](https://docs.sqlc.dev/en/latest/index.html) tool to generate type safe postgres database interface and methods.

#### `internal/usecase/microservices`
An abstract api that the usecase business logic works with. For instance this is where you do calls to external services(micro-services) hence separation of concern.
The microservice implements an interface and thus enabling mocking the interactions you would use.

## Dependency Injection
In order to remove the dependence of business logic on external packages, dependency injection is used.

For example, through the New constructor, we inject the dependency into the structure of the business logic.
This makes the business logic independent (and portable).
We can override the implementation of the interface without making changes to the `usecase` package.

```go
package usecase

type BlogUseCase struct {
	config *config.Config
	store  intfaces.Store
}

func NewBlogUseCase(store Store, config *config.Config) BlogUsecase {
	return &BlogUseCase{
		store:  store,
		config: config,
	}
}
```

### Key idea
Programmers realize the optimal architecture for an application after most of the code has been written.

> A good architecture allows decisions to be delayed to as late as possible.

### The main principle
Dependency Inversion (the same one from SOLID) is the principle of dependency inversion.
The direction of dependencies goes from the outer layer to the inner layer.
Due to this, business logic and entities remain independent from other parts of the system.

So, the application is divided into 2 layers, internal and external:
1. **Business logic** (Go standard library).
2. **Tools** (databases, servers, message brokers, any other packages and frameworks).

**The inner layer** with business logic should be clean. It should:
- Not have package imports from the outer layer.
- Use only the capabilities of the standard library.
- Make calls to the outer layer through the interface (!).

The business logic doesn't know anything about Postgres or a specific web API.
Business logic has an interface for working with an _abstract_ database or _abstract_ web API.

**The outer layer** has other limitations:
- All components of this layer are unaware of each other's existence. How to call another from one tool? Not directly, only through the inner layer of business logic.
- All calls to the inner layer are made through the interface (!).
- Data is transferred in a format that is convenient for business logic (`internal/entity`).

For example, you need to access the database from HTTP (controller).
Both HTTP and database are in the outer layer, which means they know nothing about each other.
The communication between them is carried out through `usecase` (business logic):

```
    HTTP > usecase
           usecase > repository (Postgres)
           usecase < repository (Postgres)
    HTTP < usecase
```
The symbols > and < show the intersection of layer boundaries through Interfaces.
The same is shown in the picture:

![Example](docs/img/example-http-db.png)

### Layers
![Example](docs/img/layers-2.png)

### Common terms

- **Use Cases** is business logic located in `internal/usecase`.

The layer with which business logic directly interacts is usually called the _infrastructure_ layer.
These can be repositories `internal/usecase/repo`, external webapi `internal/usecase/webapi`, any pkg, and other microservices.
In the template, the _infrastructure_ packages are located inside `internal/usecase`.

You can choose how to call the entry points as you wish. The options are:
- controller (in our case)
- delivery
- transport
- gateways
- entrypoints
- primary
- input

### Additional layers
The classic version of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) was designed for building large monolithic applications and has 4 layers.

In the original version, the outer layer is divided into two more, which also have an inversion of dependencies
to each other (directed inward) and communicate through interfaces.

The inner layer is also divided into two (with separation of interfaces), in the case of complex logic.

_______________________________

## Useful links
- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/ru/)