FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ws-server .

EXPOSE 8089

CMD ["./ws-server"]
