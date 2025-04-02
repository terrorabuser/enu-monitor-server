FROM golang:1.23.0-alpine AS build

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod tidy
RUN go build -o /app ./cmd/main.go

CMD [ "/app" ]

