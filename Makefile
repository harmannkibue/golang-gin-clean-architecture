#DB_URL=postgresql://CHURPY:CTwxdBKzThCJlzR7p169@va-gateway-prod-db.c1kdbemqpagi.us-east-1.rds.amazonaws.com:5432/VA_PROD_GATEWAY?sslmode=disable

# Adds files to staging the commits the changes and finally push to remote repository
gitPush:
	./git_push.sh

composeUp:
	docker-compose  up --build

composeDown:
	docker-compose  down --remove-orphans

# for GRPC coming soon
swagInit:
	swag init --parseInternal -g  internal/controller/http/v1/router.go

# Starts server with incorporated air tool for hot reloads
runserver:
	air -c .air.conf

# Run database schema migrations
migrateUp:
	go run -tags migrate ./cmd/app

## Rollback the changes in schema made by migrate up command
#migrateDown:
#	migrate -path ./migrations -database "$(DB_URL)" -verbose down
#
## creating migrations file in a sequential order e.g 000001_name_of_migration.up.sql, 000002_name_of_migration.up.sql ....n
#migrateCreate:
#	migrate create -ext sql -dir ./migrations -seq init_schema
#
## Goes to the specific version of the migrations. e.g version 1
#migrateGoTo:
#	migrate -path db/migration -database "$(DB_URL)" -verbose goto 1
#
#migrateDrop:
#	migrate -path ./migrations -database "$(DB_URL)" -verbose drop

# Generating mock for database -.
mockGen:
	mockgen -package mockDb -destination internal/usecase/mock/store.go github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repositories Store

# Initializing an sqlc instance -.
sqlcInit:
	cd internal && sqlc init

# For checking if there any typos in sql schemas defination -.
sqlcCompile:
	cd internal && sqlc compile

# For generating sqlc -.
sqlcGenerate:
	cd internal && sqlc generate

# For generating database schema documentations.
dbDocs:
	dbdocs build docs/db.dbml

# Phony avoids conflicts for the file named as the main command
.PHONY: gitPush composeUp composeDown composeBuildAdminApp swagInit runserver migrateUp migrateDown migrateCreate
.PHONY: migrateGoTo migrateDrop sqlcInit sqlcCompile sqlcGenerate dbDocs
