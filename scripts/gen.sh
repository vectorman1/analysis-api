go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

cd ..

mkdir -p api/swagger/v1

mkdir -p generated

# Instrument Service
mkdir -p generated/instrument_service
protoc -I . --proto_path=domain/instrument/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/instrument_service \
            --go_opt=paths=source_relative \
            --go-grpc_out=generated/instrument_service \
            --go-grpc_opt=paths=source_relative \
            instrument_service.proto

protoc -I . --proto_path=domain/instrument/proto/v1 \
            --proto_path=third_party \
            --openapiv2_out ./api/swagger/v1 \
            --openapiv2_opt logtostderr=true \
            instrument_service.proto

protoc -I . --proto_path=domain/instrument/proto/v1 \
            --proto_path=third_party \
            --grpc-gateway_out ./generated/instrument_service \
            --grpc-gateway_opt paths=source_relative \
            --grpc-gateway_opt generate_unbound_methods=true \
            instrument_service.proto

# User Service
mkdir -p generated/user_service
protoc -I . --proto_path=domain/user/proto/v1 \
            --proto_path=third_party \
            --go_out=generated/user_service \
            --go_opt=paths=source_relative \
            --go-grpc_out=generated/user_service \
            --go-grpc_opt=paths=source_relative \
            user_service.proto

protoc -I . --proto_path=domain/user/proto/v1 \
            --proto_path=third_party \
            --openapiv2_out ./api/swagger/v1 \
            --openapiv2_opt logtostderr=true \
             user_service.proto

protoc -I . --proto_path=domain/user/proto/v1 \
            --proto_path=third_party \
            --grpc-gateway_out ./generated/user_service \
            --grpc-gateway_opt paths=source_relative \
            --grpc-gateway_opt generate_unbound_methods=true \
             user_service.proto
