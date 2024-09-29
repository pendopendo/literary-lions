FROM golang:1.21

WORKDIR /literary-lions
COPY . .

RUN go mod download

RUN go build -o main .

ENTRYPOINT ["/literary-lions/main"]

LABEL name="literary-lions"