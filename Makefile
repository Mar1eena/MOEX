name = trb

# Переменные для подключения к базе данных
DB_USER=admin
DB_PASSWORD=1984
DB_HOST=postgres
DB_PORT=5432
DB_NAME=TrB_DB
SCALE_WS=1
SCALE_GRPC=1
SCALE_gRPC_receiving_SERVER=1
SCALE_gRPC_PUBLISHING_SERVER=1

# Строка подключения
DB_CONNECTION_STRING=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

build:
	docker-compose --project-name=${name} build

up:
	docker-compose --project-name=${name} up -d --scale server-ws=$(SCALE_WS)

upd:
	docker-compose --project-name=${name} up --build -d
down:
	docker-compose --project-name=${name} down

rungrpc:
	cd ./cmd/server-grpc-publishing && go run main.go 

runws:
	cd ./cmd/server_ws && go run main.go 

proto-gen-server-grpc-publishing:
	protoc -I internal/pkg/services/server-grpc-publishing/proto \
  	internal/pkg/services/server-grpc-publishing/proto/*.proto \
	--go_out=./internal/pkg/services/server-grpc-publishing/proto --go_opt=paths=source_relative \
	--go-grpc_out=./internal/pkg/services/server-grpc-publishing/proto  --go-grpc_opt=paths=source_relative \

proto-gen-server-grpc-receiving:
	protoc -I internal/pkg/services/server-grpc-receiving/client/tinkoff-invest/proto \
  	internal/pkg/services/server-grpc-receiving/client/tinkoff-invest/proto/*.proto \
	--go_out=./internal/pkg/services/server-grpc-receiving/client/tinkoff-invest/proto --go_opt=paths=source_relative \
	--go-grpc_out=./internal/pkg/services/server-grpc-receiving/client/tinkoff-invest/proto  --go-grpc_opt=paths=source_relative \

	