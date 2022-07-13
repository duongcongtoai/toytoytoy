
echo "hello";
sleep 3000s;
echo "hello";

CGO_ENABLED=0 go test ./cmd/test/...  -v -count 1 --coverprofile=./tests/report/integration.out \
--coverpkg=./internal/... \
&& go tool cover -html ./tests/report/integration.out -o ./tests/report/integration.html

# CGO_ENABLED=0 go test ./internal/tests/storage/... -v --coverprofile=./tests/report/storage.out \
# --coverpkg=./internal/storages/psql/... \
# && go tool cover -html ./tests/report/storage.out -o ./tests/report/storage.html

# CGO_ENABLED=0 go test ./internal/tests/integration/... -v --coverprofile=./tests/report/integration.out \
# --coverpkg=./internal/usecase/...,./internal/storages/...,./internal/transport/... \
# && go tool cover -html ./tests/report/integration.out -o ./tests/report/integration.html