FROM golang:1.18-alpine as builder
WORKDIR /src/go

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./internal ./internal
COPY ./sqlc ./sqlc
COPY ./cmd ./cmd


RUN CGO_ENABLED=0 GOOS=linux go build -a -o server ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o migrate ./cmd/migration

CMD ["./server"]
