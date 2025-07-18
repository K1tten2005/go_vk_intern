FROM golang:1.24.2-alpine AS builder

COPY . /github.com/K1tten2005/go_vk_intern
WORKDIR /github.com/K1tten2005/go_vk_intern

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main/main.go
RUN go clean --modcache

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/K1tten2005/go_vk_intern/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080

ENTRYPOINT ["./.bin"]