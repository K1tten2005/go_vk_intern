FROM golang:1.24.2-alpine AS builder

COPY . /github.com/K1tten2005/go_vk_intern
WORKDIR /github.com/K1tten2005/go_vk_intern

RUN apk add --no-cache ca-certificates
RUN go mod download
RUN mkdir -p ./logs
RUN touch ./logs/main.log
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go
RUN go clean --modcache

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/K1tten2005/go_vk_intern/.bin .
COPY --from=builder /github.com/K1tten2005/go_vk_intern/logs ./logs  
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080
ENTRYPOINT ["./.bin"]
