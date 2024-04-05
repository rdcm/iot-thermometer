FROM golang:1.22.2 AS build-stage
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer


FROM golang:1.22.2-alpine3.19
WORKDIR /
COPY --from=build-stage /app/consumer   /consumer
COPY --from=build-stage /app/config.yml /config.yml

ENTRYPOINT ["/consumer"]