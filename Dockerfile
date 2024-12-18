FROM golang:1.13.6-alpine

ENV APP_NAME ip-analyser
ENV APP_PORT 8080

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go get ./
RUN go build -o ${APP_NAME}

CMD ./${APP_NAME}

EXPOSE ${APP_PORT}