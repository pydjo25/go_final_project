FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o my-todo main.go

ARG PORT=7540
EXPOSE $PORT

ENV TODO_PORT=$PORT
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

VOLUME /data

CMD ["./my-todo"]