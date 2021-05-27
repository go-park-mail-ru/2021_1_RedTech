FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main_sub cmd/subscription/main.go


FROM alpine

WORKDIR /app

COPY --from=builder app/main_sub .
EXPOSE 8084

CMD sleep 10 && ./main_sub
