FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ARG TOKEN=not_set
ARG VER_TOKEN=not_set
ARG MONGO=not_set
ENV TOKEN=$TOKEN
ENV MONGO=$MONGO
ENV VER_TOKEN=$VER_TOKEN

EXPOSE 8080

CMD ["app"]
