FROM golang:1.22.2 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer


FROM golang:1.22.2-alpine3.19
WORKDIR /app
COPY --from=build /build/consumer   /app/consumer
COPY --from=build /build/config.yml /app/config.yml
ENTRYPOINT ["/app/consumer"]