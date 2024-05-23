FROM golang:1.21.1

WORKDIR /app

COPY . .

RUN go mod download

COPY *.go ./

RUN go build -o /main

EXPOSE 5000

CMD [ "/main" ]
