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
RUN MAIN_FILE=${MAIN_FILE} && \
    go build -o app ${MAIN_FILE}

FROM alpine:latest AS runtime
RUN apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/app /app/app

CMD ["/app/app", "server"]
