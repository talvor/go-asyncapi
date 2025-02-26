db_login:
	psql ${DATABASE_URL}

db_test_login:
	psql ${DATABASE_TEST_URL}

db_create_migration:
	migrate create -ext sql -dir migrations -seq $(name)

db_migrate:
	migrate -database ${DATABASE_URL} -path migrations up

stack_up:
	podman compose up -d

stack_down:
	podman compose down

run_unit_tests:
	TESTCONTAINERS_RYUK_DISABLED=true ginkgo run -v ./...

start_apiserver:
	go run cmd/apiserver/main.go

terraform_apply:
	terraform -chdir=terraform apply -auto-approve
