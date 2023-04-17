include .env.example
export

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# Phony avoids conflicts for the file named as the main command
.PHONY: help
# Adds files to staging the commits the changes and finally push to remote repository

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

gitPush: ### Command to aid in pushing code to the current branch. Runs `git add . && git commit && git push
	./git_push.sh
.PHONY: gitPush

composeUp: ## Command to build and run docker containers
	docker-compose  up --build
.PHONY: composeUp

composeDown: ### Command to stop docker containers
	docker-compose  down --remove-orphans
.PHONY: composeDown

swagInit: ### Command used to initialize swagger api documentation.Uses swaggo and gin swagger packages
	swag init --parseInternal -g  internal/controller/http/v1/router.go
.PHONY: swagInit

runserver: ### Starts server with incorporated air tool for hot reloads
	air -c .air.conf
.PHONY: runserver

migrateUp: ### Command to run database schema migrations.The package used to run migrations is golang migrate
	go run -tags migrate ./cmd/app
.PHONY: migrateUp

migrateDown: ### command to rollback the changes in schema made by migrate up command
	migrate -path ./migrations -database "$(DB_URL)" -verbose down
.PHONY: migrateDown

migrateCreate: ### Command for creating migrations file in a sequential order e.g 000001_name_of_migration.up.sql, 000002_name_of_migration.up.sql ....n
	migrate create -ext sql -dir ./migrations -seq init_schema
.PHONY: migrateCreate

migrateGoTo: ### Goes to the specific version of the migrations. e.g version 1
	migrate -path db/migration -database "$(DB_URL)" -verbose goto 1
.PHONY: migrateGoTo

migrateDrop: ### Command used to drop the database migrations
	migrate -path ./migrations -database "$(DB_URL)" -verbose drop
.PHONY: migrateDrop

mockGen: ### Generating mock for database mock testing.The package used is goMock
	mockgen -package mockDb -destination internal/usecase/mock/store.go github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories Store
.PHONY: mockGen

sqlcInit: ### Command used to initialize an database query repository store.Package used is SQLC.
	cd internal && sqlc init
.PHONY: sqlcInit

sqlcCompile: ### Command used for checking if there any typos in sql schemas definition
	cd internal && sqlc compile
.PHONY: sqlcCompile

sqlcGenerate: ### Command used to generate database query repository store.Package used is SQLC -.
	cd internal && sqlc generate
.PHONY: sqlcGenerate

dbDocs: ### Command used for documenting generating database schema documentations.
	dbdocs build docs/db.dbml
.PHONY: dbDocs

mockeryGenerateBlogUsecase: ### generates testing mocks using mockery tool
	mockery --dir=internal/entity/intfaces --name=BlogUsecase --filename=blog.go --output=internal/entity/mocks --outpkg=mocks
.PHONY: mockeryGenerateBlogUsecase