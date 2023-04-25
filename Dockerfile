FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o kv_ts ./cmd/


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/kv_ts kv_ts


EXPOSE 8080

ENTRYPOINT ["/kv_ts"]
