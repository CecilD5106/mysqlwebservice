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

RUN mkdir /app

COPY --from=builder /dist/mysqlweb /app/
COPY ./form/Edit.tmpl /app/form/Edit.tmpl
COPY ./form/Footer.tmpl /app/form/Footer.tmpl
COPY ./form/Header.tmpl /app/form/Header.tmpl
COPY ./form/Index.tmpl /app/form/Index.tmpl
COPY ./form/Menu.tmpl /app/form/Menu.tmpl
COPY ./form/New.tmpl /app/form/New.tmpl
COPY ./form/Show.tmpl /app/form/Show.tmpl

EXPOSE 8080

ENTRYPOINT ["/app/mysqlweb"]