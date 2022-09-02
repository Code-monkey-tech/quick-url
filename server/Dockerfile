FROM golang:1.18-alpine AS build

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go build -o bin/shrty ./cmd/main.go

EXPOSE ${PORT}
CMD ["bin/shrty"]