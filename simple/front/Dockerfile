FROM golang:1.11.5-alpine3.9
LABEL version="006"
LABEL vendor="Sounding"
EXPOSE 8888
WORKDIR /go/src/craque/
COPY . .
RUN apk add --no-cache git
RUN go get github.com/rs/zerolog
RUN go get github.com/prometheus/client_golang/prometheus
RUN go build -o /bin/craque
ENTRYPOINT ["/bin/craque"]
