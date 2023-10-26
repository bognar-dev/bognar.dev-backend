FROM golang:1.19.5-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./out/dist
CMD ["./out/dist"]