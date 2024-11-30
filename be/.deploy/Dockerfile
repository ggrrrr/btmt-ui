# ARG SERVICE
# FROM golang:1.21 AS builder
# FROM golang:1.21-alpine AS builder
FROM golang:1.22-alpine AS builder

# RUN apt update && apt install telnet curl
WORKDIR /app

# COPY go.* ./

COPY  . ./
RUN go mod tidy
RUN go mod download

ARG MAIN_FILE
ARG GIT_HASH
ARG VER_PACKAGE=github.com/ggrrrr/btmt-ui/be/common/buildversion

RUN go build  -ldflags "-X ${VER_PACKAGE}.version=${GIT_HASH}" \
    -o app ${MAIN_FILE}

FROM alpine:latest AS runtime
ARG GIT_HASH
ARG MAIN_FILE

RUN apk add --no-cache curl
WORKDIR /app

RUN addgroup -g 1000 -S app && \
    adduser -u 1000 -h /app -G app -S app

COPY --from=builder --chown=app:app /app/app /app/app

# COPY --from=builder /app/config.yaml /app/config.yaml
# COPY --from=builder /app/jwt.key /app/jwt.key
# COPY --from=builder /app/jwt.crt /app/jwt.crt

CMD ["/app/app", "server"]

## Some labels
LABEL org.opencontainers.image.title="${MAIN_FILE}" \
    org.opencontainers.image.revision="${GIT_HASH}" \
    org.opencontainers.image.vendor="btmt.io"
