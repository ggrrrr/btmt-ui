# ARG SERVICE
# FROM golang:1.21 AS builder
# FROM golang:1.21-alpine AS builder
FROM golang:1.21-alpine AS builder

# RUN apt update && apt install telnet curl
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY  . ./

ARG MAIN_FILE

# RUN echo $SERVICE
# RUN date -u +"%Y-%m-%dT%H:%M:%SZ" > be/common/ver/ts.txt 
RUN MAIN_FILE=${MAIN_FILE} && \
    go build -o app ${MAIN_FILE}

FROM alpine:latest AS runtime
RUN apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/app /app/app
# COPY --from=builder /app/config.yaml /app/config.yaml
# COPY --from=builder /app/jwt.key /app/jwt.key
# COPY --from=builder /app/jwt.crt /app/jwt.crt

CMD ["/app/app", "server"]
