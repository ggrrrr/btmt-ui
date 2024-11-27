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
# RUN echo $SERVICE

# RUN MAIN_FILE=${MAIN_FILE} && \
RUN go build  -ldflags "-X ${VER_PACKAGE}.version=${GIT_HASH}" \
    -o app ${MAIN_FILE}

FROM alpine:latest AS runtime
ARG GIT_HASH

RUN apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/app /app/app

# this is used by common/system/system.go during startup

# COPY --from=builder /app/config.yaml /app/config.yaml
# COPY --from=builder /app/jwt.key /app/jwt.key
# COPY --from=builder /app/jwt.crt /app/jwt.crt

CMD ["/app/app", "server"]
