run:
	go run ./cmd/api

network:
	docker network create simpleblog-network

postgres:
	docker run --name simple_blog_db --network simpleblog-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it simple_blog_db createdb --username=root --owner=root simpleblog

dropdb:
	docker exec -it simple_blog_db dropdb simpleblog

genmi:
	migrate create -ext sql -dir db/migration -seq init_schema

mgu:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleblog?sslmode=disable" -verbose up

mgd:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleblog?sslmode=disable" -verbose down

sqlc:
	sqlc generate
