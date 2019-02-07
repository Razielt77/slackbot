FROM golang:1.8

RUN cd /usr/local/bin && \
    wget -qO- https://github.com/codefresh-io/cli/releases/download/v0.13.2/codefresh-v0.13.2-linux-x64.tar.gz\
     | tar xvz

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ARG TOKEN=not_set
ENV TOKEN=$TOKEN

EXPOSE 8080

CMD ["app"]
