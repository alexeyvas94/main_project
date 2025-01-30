LOCAL_BIN:=$(CURDIR)/bin #CURDIR — это встроенная переменная, которая содержит путь к текущей директории

install-deps_proto:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-car:
	mkdir -p pkg/car #Создаёт папку pkg/car, если она ещё не существует.
	protoc --proto_path=. \
	--go_out=pkg/car --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/car --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	car.proto

	#protoc --proto_path=.
	#Указывает путь, где искать Protobuf-файлы. В данном случае в корне проекта.
	#--go_out=pkg/car
	#Указывает, куда сгенерировать Go-код для Protobuf-сообщений.
	#В данном случае файлы будут созданы в pkg/car.
	#--go_opt=paths=source_relative
	#Заставляет protoc генерировать Go-код с относительными путями для импортов, что упрощает работу в многомодульных проектах.
	#--plugin=protoc-gen-go=bin/protoc-gen-go
	#Указывает явный путь к плагину protoc-gen-go (используется для генерации структур сообщений .pb.go).
	#--go-grpc_out=pkg/car
	#Указывает, куда сгенерировать Go-код для gRPC-клиента и сервера.
	#gRPC-код будет создан в той же директории pkg/car.
	#--go-grpc_opt=paths=source_relative
	#Аналогично --go_opt=paths=source_relative, но применяется к gRPC-коду.
	#--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc
	#Указывает путь к плагину protoc-gen-go-grpc, который генерирует код для gRPC.
	#car.proto
	#Указывает путь к конкретному .proto файлу, который будет использоваться для генерации.
include .env
LOCAL_BIN:=$(CURDIR)/bin
install-deps:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
