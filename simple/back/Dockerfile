FROM golang:1.11.5-alpine3.9
LABEL version="008"
LABEL vendor="Sounding"
EXPOSE 9999
WORKDIR /go/src/bacque/
COPY . .
RUN apk add --no-cache git
RUN go get github.com/rs/zerolog
RUN go get github.com/prometheus/client_golang/prometheus
RUN go build -o /bin/bacque
ENTRYPOINT ["/bin/bacque"]
