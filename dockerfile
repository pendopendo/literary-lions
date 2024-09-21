FROM golang:1.21

WORKDIR /literary-lions-forum

COPY . .

RUN go mod download

RUN go build -o main .

ENTRYPOINT ["/literary-lions-forum/main"]

LABEL name="literary-lions-forum"