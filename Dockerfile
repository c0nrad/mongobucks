FROM golang:alpine

MAINTAINER stuart.larsen <stuart@mongodb.com>

RUN apk add --update \
    git

RUN go get "github.com/codegangsta/negroni" "github.com/gorilla/context" "gopkg.in/mgo.v2" "github.com/gorilla/sessions"
RUN go get "golang.org/x/oauth2" "golang.org/x/oauth2/google" "gopkg.in/mgo.v2/bson"

WORKDIR /go/src/github.com/c0nrad/mongobucks
COPY . ./

RUN go get
RUN go build

EXPOSE 8081
CMD ["./mongobucks"]
