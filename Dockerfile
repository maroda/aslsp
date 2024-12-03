FROM golang:alpine3.20
LABEL app="aslsp"
LABEL version="0.4.0"
LABEL vendor="Sounding"
EXPOSE 8888
EXPOSE 9999
WORKDIR /go/src/aslsp/
COPY . .
RUN go get .
RUN go build -o /bin/aslsp
ENTRYPOINT ["/bin/aslsp"]
