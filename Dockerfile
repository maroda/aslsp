FROM golang:alpine3.13
LABEL vendor="Sounding"
LABEL version="0.2.0"
EXPOSE 8888
EXPOSE 9999
WORKDIR /go/src/aslsp/
COPY . .
RUN go get .
RUN go build
