FROM golang:1.11.5-alpine3.9
LABEL version="001"
LABEL vendor="Sounding"
EXPOSE 7777
WORKDIR /go/src/inque/
COPY . .
RUN go build -o /bin/inque
ENTRYPOINT ["/bin/inque"]
