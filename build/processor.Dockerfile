FROM golang:1.22-alpine as builder
WORKDIR /src
COPY cmd/test_server/main.go ./main.go
COPY internal/models ./internal/models
COPY internal/client/queue ./internal/client/queue
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
RUN go build -o bin/test_server ./main.go

FROM alpine:edge
COPY --from=builder src/bin/test_server bin/test_server
CMD ["/bin/test_server"]
