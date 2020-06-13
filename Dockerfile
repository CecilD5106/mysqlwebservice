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
COPY ./form/Edit.tmpl /dist/form/Edit.tmpl
COPY ./form/Footer.tmpl /dist/form/Footer.tmpl
COPY ./form/Header.tmpl /dist/form/Header.tmpl
COPY ./form/Index.tmpl /dist/form/Index.tmpl
COPY ./form/Menu.tmpl /dist/form/Menu.tmpl
COPY ./form/New.tmpl /dist/form/New.tmpl
COPY ./form/Show.tmpl /dist/form/Show.tmpl

EXPOSE 8080

CMD ["/dist/mysqlweb"]