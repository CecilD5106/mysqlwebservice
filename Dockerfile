FROM golang:alpine as builder

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

FROM alpine

COPY --from=builder /dist/mysqlweb /
COPY ./form/Edit.tmpl /form/Edit.tmpl
COPY ./form/Footer.tmpl /form/Footer.tmpl
COPY ./form/Header.tmpl /form/Header.tmpl
COPY ./form/Index.tmpl /form/Index.tmpl
COPY ./form/Menu.tmpl /form/Menu.tmpl
COPY ./form/New.tmpl /form/New.tmpl
COPY ./form/Show.tmpl /form/Show.tmpl

EXPOSE 8080

ENTRYPOINT ["/mysqlweb"]