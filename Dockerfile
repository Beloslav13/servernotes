FROM golang:1.18-buster as build

# set work directory
WORKDIR /app

COPY . /app

EXPOSE 8080

RUN go mod download
RUN go build -o /go/notes /app/cmd/app/main.go

# ????
CMD ["./notes"]