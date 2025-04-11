FROM golang:1.24-alpine as builder

WORKDIR /builder

COPY . .

RUN go help modules
RUN go build --gcflags="all=-N -l" -o ./main ./cmd/main.go

FROM alpine:3.21

WORKDIR /app

EXPOSE 3000

COPY --from=builder /builder/main .
COPY --from=builder /builder/configs/ ./configs 

ENTRYPOINT ./main -config=configs/docker.yaml