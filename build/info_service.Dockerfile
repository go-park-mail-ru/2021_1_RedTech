FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/info/main.go


FROM alpine

WORKDIR /app

COPY --from=builder app/main .
EXPOSE 8081

CMD sleep 10 && ./main