
gen-mock:
	mockery --dir internal/services --all --output ./mock/services --outpkg services
	mockery --dir internal/common --all --output ./mock/common --outpkg common
gen-sqlc:
	sqlc generate -f sqlc/sqlc.yaml

test-integration-local:
	CGO_ENABLED=0 go test ./cmd/test/... --config config.local.yaml --serverhost localhost -v -count 1 \
	--coverprofile=./tests/report/integration.out \
	--coverpkg=./internal/... \
	&& go tool cover -html ./tests/report/integration.out -o ./tests/report/integration.html
