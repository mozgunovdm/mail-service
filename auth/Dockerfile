FROM golang:1.17-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o application ./cmd/auth-service/

FROM alpine:3.15.4
COPY --from=builder /app/application /app/application
COPY --from=builder /app/user_config.json /user_config.json
COPY --from=builder /app/config.yml /config.yml
CMD ["/app/application"]