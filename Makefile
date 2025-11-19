include .env.example
export

DB_URL = "postgresql://postgres:password@postgres:5432/BLOG_DB?sslmode=disable"

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

swagInit: ### Command used to initialize swagger api documentation.Uses swaggo and gin swagger packages.Run go get -u github.com/swaggo/swag/cmd/swag@v1.6.7 command first
	swag init --parseInternal --parseInternal -g  internal/controller/http/v1/router.go
.PHONY: swagInit

runserver: ### Starts server with incorporated air tool for hot reloads
	air -c .air.conf
.PHONY: runserver

dbDocs: ### Command used for documenting generating database schema documentations.
	dbdocs build docs/db.dbml
.PHONY: dbDocs

mockeryGenerateBlogUsecase: ### generates testing mocks using mockery tool
	mockery --dir=internal/entity/intfaces --name=BlogUsecase --filename=blog.go --output=internal/entity/mocks --outpkg=mocks
.PHONY: mockeryGenerateBlogUsecase

testWithCoverProfile: ### Used to run tests with coverage and display the output.Scans all the files and runs the tests if available
	go test ./... -coverprofile=cover.out
.PHONY: goTestCoverProfile
