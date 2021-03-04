go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

mkdir -p api/swagger/v1

mkdir -p generated

mkdir -p generated/proto_models
protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/proto_models \
            --go_opt=paths=source_relative \
            models.proto


mkdir -p generated/symbol_service
# Symbols Service
protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/symbol_service \
            --go_opt=paths=source_relative \
            --go-grpc_out=generated/symbol_service \
            --go-grpc_opt=paths=source_relative \
            symbol_service.proto

protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --openapiv2_out ./api/swagger/v1 \
            --openapiv2_opt logtostderr=true \
            symbol_service.proto

protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --grpc-gateway_out ./generated/symbol_service \
            --grpc-gateway_opt paths=source_relative \
            --grpc-gateway_opt generate_unbound_methods=true \
            symbol_service.proto

mkdir -p generated/user_service
protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/user_service \
            --go_opt=paths=source_relative \
            --go-grpc_out=generated/user_service \
            --go-grpc_opt=paths=source_relative \
            user_service.proto

protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --openapiv2_out ./api/swagger/v1 \
            --openapiv2_opt logtostderr=true \
             user_service.proto

protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --grpc-gateway_out ./generated/user_service \
            --grpc-gateway_opt paths=source_relative \
            --grpc-gateway_opt generate_unbound_methods=true \
             user_service.proto

mkdir -p generated/historical_service
protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/historical_service \
            --go_opt=paths=source_relative \
            --go-grpc_out=generated/historical_service \
            --go-grpc_opt=paths=source_relative \
            historical_service.proto

protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --openapiv2_out ./api/swagger/v1 \
            --openapiv2_opt logtostderr=true \
             historical_service.proto

protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --grpc-gateway_out ./generated/historical_service \
            --grpc-gateway_opt paths=source_relative \
            --grpc-gateway_opt generate_unbound_methods=true \
            historical_service.proto

# Worker services

mkdir -p generated/worker_symbol_service
protoc -I . --proto_path=api/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/worker_symbol_service \
            --go_opt=paths=source_relative \
            --go-grpc_out=generated/worker_symbol_service \
            --go-grpc_opt=paths=source_relative \
            worker_symbol_service.proto