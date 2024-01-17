.PHONY: migrate
migrate:
	cd sql/schema && goose postgres postgres://postgres:admin@localhost:5432/rssfeed up

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate


.PHONY: t-watch
t-watch:
	templ generate -watch

.PHONY: serve
serve:
	air

.PHONY: build 
build:
	go build
