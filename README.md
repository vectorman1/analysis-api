# analysis-api

<hr>

gGRPC service + HTTP Gateway API for Analysis

Deployment through Github Actions + webooks to docker host.

Use `./scripts/gen.sh` to generate Swagger and gRPC service definitions.

Protobuf definitions are found in the respective domain's services folder under `./domain/{domain-name}/proto/v1/`

Service specs can be found then in `./api/swagger/v1/{service-name}.swagger.json` - they can then be imported in Postman or similar.

Explore the swagger spec for available endpoints.
