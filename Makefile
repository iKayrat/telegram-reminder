.PHONY: run psqlconn create drop migrateup migratedown sqlc
run:
	go run main.go
#------------------
# createdb:
# 	 winpty docker exec -it postgres13 createdb --username=root ecom_kg

# dropdb:
# 	winpty docker exec -it postgres13  dropdb --username=root ecom_kg

# createmigration:
# 	migrate create -ext sql -dir model/migration -seq name_table

migrateup:
	migrate -path migrations -database "postgresql://root:kaak@localhost:5432/tgclients?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:kaak@localhost:5432/tgclients?sslmode=disable" -verbose down

#connnect to docker conntainer postgres13 dbname ecom_kg
psqlconn:
	docker exec -it telegramdb psql -U root tgclients

sqlc:
	sqlc generate

postgres:
	docker run --name telegramdb -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=kaak -d postgres:alpine
create:
	docker exec -it telegramdb createdb --username=root --owner=root tgclients
drop:
	docker exec -it telegramdb dropdb --username=root tgclients