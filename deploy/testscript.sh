CGO_ENABLED=0 go test ./cmd/test/...  -v -count 3 --coverprofile=./tests/report/integration.out \
--coverpkg=./internal/... \
&& go tool cover -html ./tests/report/integration.out -o ./tests/report/integration.html

CGO_ENABLED=0 go test ./internal/services/... -v -count 3 --coverprofile=./tests/report/unittest.out \
--coverpkg=./internal/services/... \
&& go tool cover -html ./tests/report/unittest.out -o ./tests/report/unittest.html

