FROM golang:1.13 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o server ./cmd/gograph
RUN CGO_ENABLED=0 GOOS=linux make createdb

FROM alpine:3
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/bin/createdb /app/createdb

CMD ["/app/createdb", "--database", "gs://matsf-data/movies.db"]
