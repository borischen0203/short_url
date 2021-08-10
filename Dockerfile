FROM golang:1.16-alpine
WORKDIR /app
ADD . /app

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN cd /app && go build -o main

EXPOSE 8000

CMD [ "/app/main" ]
