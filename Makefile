db-init:
	psql -c 'CREATE DATABASE "balam"' -U $(user)

db-init-kube:
	psql -c 'CREATE DATABASE "balam"' -h $(host) -p $(port) -U $(user)

migrationup:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/balam?sslmode=disable" -verbose up

migrationdown:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/balam?sslmode=disable" -verbose down

migrationup1:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/balam?sslmode=disable" -verbose up 1

migrationdown1:
	migrate -path db/migrations -database "postgres://$(user):$(password)@$(host):$(port)/balam?sslmode=disable" -verbose down 1

.PHONY: db-init db-init-kube migrationup migrationdown migrationup1 migrationdown1