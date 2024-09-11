FROM golang:1.23.1-bookworm AS developement

RUN mkdir /app
WORKDIR /app

RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8080
CMD ["air","-c",".air.toml"]

