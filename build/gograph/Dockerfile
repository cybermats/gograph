FROM golang:1.13 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o server ./cmd/gograph
RUN CGO_ENABLED=0 GOOS=linux make gograph

FROM alpine:3
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/bin/gograph /app/gograph
COPY --from=builder /app/web /web

CMD ["/app/gograph", "--database", "gs://matsf-data/movies.db", "--web-dir", "/app/web"]
