FROM golang:1.11.3-alpine3.8
LABEL net.craque.version="003"
LABEL vendor="Sounding"
EXPOSE 8888
WORKDIR /go/src/craque/
COPY . .
RUN apk add --no-cache git
RUN go get github.com/rs/zerolog
RUN go build -o /bin/craque
ENTRYPOINT ["/bin/craque"]
