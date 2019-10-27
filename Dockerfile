FROM golang:1.13 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o server ./cmd/gograph
RUN go run ./cmd/gograph --filename database.gz --create

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server
COPY --from=builder /app/database.gz /database.gz
COPY --from=builder /app/web /web

CMD ["/server", "--filename", "/database.gz", "--web-dir", "/web"]
