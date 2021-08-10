FROM golang:1.16-alpine
# RUN apk add --no-cache git
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /short_url

EXPOSE 8000

CMD [ "/short_url" ]
