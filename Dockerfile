FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o mysqlweb .

WORKDIR /dist

RUN cp /build/mysqlweb .
RUN cp /build/form/Edit.tmpl .
RUN cp /build/form/Footer.tmpl .
RUN cp /build/form/Header.tmpl .
RUN cp /build/form/Index.tmpl .
RUN cp /build/form/Menu.tmpl .
RUN cp /build/form/New.tmpl .
RUN cp /build/form/Show.tmpl .

EXPOSE 8080

CMD ["/dist/mysqlweb"]