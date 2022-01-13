postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root kez_db

dropdb:
	docker exec -it postgres14 dropdb kez_db

migrateup:
	migrate -path ./schema -database 'postgresql://root:123456@0.0.0.0:5432/kez_db?sslmode=disable' -verbose up

migratedown:
	migrate -path ./schema -database 'postgresql://root:123456@0.0.0.0:5432/kez_db?sslmode=disable' -verbose down


.PHONY: postgres createdb dropdb migrateup migratedown
