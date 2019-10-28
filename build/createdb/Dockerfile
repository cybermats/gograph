FROM golang:1.13 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o server ./cmd/gograph
#RUN go run ./cmd/gograph --filename database.gz --create
RUN make createdb

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/ /
COPY --from=builder /app/web /web

CMD ["/createdb", "--filename", "gs://matsf-data/movies.db"]
