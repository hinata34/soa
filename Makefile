
.PHONY: generate_proto
generate_proto:
	protoc --proto_path=./internal/userservice/proto --go_out=./internal/userservice/proto/pb --go_opt=paths=source_relative --go-grpc_out=./internal/userservice/proto/pb --go-grpc_opt=paths=source_relative userservice.proto

.PHONY: psql_cli
psql_cli: 
	psql -h localhost -p 5432 -d test -U test

migration_up:
	goose -dir ./migrations postgres "postgresql://test:test@127.0.0.1:5432/test?sslmode=disable" up

migration_down:
	goose -dir ./migrations postgres "postgresql://test:test@127.0.0.1:5432/test?sslmode=disable" down

generate_openapi:
	oapi-codegen -config ./internal/apigateway/swagger/server.cfg.yaml -o ./internal/apigateway/swagger/server.gen.go internal/apigateway/swagger/openapi.yaml

.PHONY: start
start:
	docker compose build && docker compose up