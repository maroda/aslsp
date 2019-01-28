FROM golang:1.11.3-alpine3.8
LABEL net.bacque.version="003"
LABEL vendor="Sounding"
EXPOSE 9999
WORKDIR /go/src/bacque/
COPY . .
RUN apk add --no-cache git
RUN go get github.com/rs/zerolog
RUN go build -o /bin/bacque
ENTRYPOINT ["/bin/bacque"]
