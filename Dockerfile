FROM golang:alpine3.13
LABEL app="aslsp"
LABEL version="0.2.4"
LABEL vendor="Sounding"
EXPOSE 8888
EXPOSE 9999
WORKDIR /go/src/aslsp/
COPY . .
RUN go get .
RUN go build -o /bin/aslsp
ENTRYPOINT ["/bin/aslsp"]
