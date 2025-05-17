export
include .env

LOCAL_BIN:=$(CURDIR)/bin
PROTO_DIR := api/auth_api
PROTO_FILE := $(PROTO_DIR)/auth.proto
OUT_DIR := pkg/w1

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make vendor-proto
	make generate-note-api
	make generate-auth-api
	make generate-access-api

generate-note-api:
	mkdir -p $(OUT_DIR)
	protoc --proto_path=$(PROTO_DIR) \
	--proto_path=vendor.protogen \
	--proto_path=vendor.protogen/google \
	--proto_path=vendor.protogen/validate \
	--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:$(OUT_DIR) --validate_opt=paths=source_relative \
    --plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
    --plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	$(PROTO_FILE)

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth_v1/authv1.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/access_v1/access.proto

local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

test:
	go clean -testcache
	go test \
      ./internal/api/tests \
      ./service/auth/tests \
      -covermode=count \
      -coverpkg=./internal/api/auth,./service/auth \
      -coverprofile=cover.out \
      -v
test-coverage:
	go clean -testcache
	go test \
		./internal/api/tests \
		./service/auth/tests \
		-covermode=count \
		-coverpkg=./internal/api/auth,./service/auth \
		-coverprofile=coverage.tmp.out \
		-count=5 \
		-v
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
        			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
        			mkdir -p  vendor.protogen/google/ &&\
        			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
        			rm -rf vendor.protogen/googleapis ;\
        		fi

gen-cert:
	mkdir -p certificates
	# Генерация закрытого ключа для центра сертификации (CA)
	openssl genrsa -out certificates/ca.key 4096
	# Создание самоподписанного сертификата CA на основе ключа (срок действия 365 дней)
	openssl req -new -x509 -key certificates/ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out certificates/ca.cert
	# Генерация закрытого ключа для gRPC-сервиса
	openssl genrsa -out certificates/service.key 4096
	# Создание запроса на сертификат (CSR) для сервиса с использованием конфигурационного файла
	openssl req -new -key certificates/service.key -out certificates/service.csr -config certificate.conf
	# Подписание CSR с помощью CA, чтобы выпустить действительный сертификат для сервиса
	# Используется конфигурационный файл для указания расширений сертификата
	openssl x509 -req -in certificates/service.csr -CA certificates/ca.cert -CAkey certificates/ca.key -CAcreateserial \
    		-out certificates/service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext

grpc-load-test:
	ghz \
		--proto api/auth_api/auth.proto \
		-i api/auth_api,vendor.protogen,vendor.protogen/google,vendor.protogen/validate \
		--call auth_api.UserAPI.Get \
		--data '{"id": 2}' \
		--rps 100 \
		--total 3000 \
		--cacert certificates/ca.cert \
        --cert certificates/service.pem \
        --key certificates/service.key \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/auth_api/auth.proto \
		-i api/auth_api,vendor.protogen,vendor.protogen/google,vendor.protogen/validate \
		--call auth_api.UserAPI.Get \
		--data '{"id": 0}' \
		--rps 100 \
		--total 3000 \
		--cacert certificates/ca.cert \
        --cert certificates/service.pem \
        --key certificates/service.key \
		localhost:50051