FROM golang:1.16 as build

WORKDIR /build

COPY . .
RUN go get .
RUN CGO_ENABLED=0 go build -o mysql-console

FROM scratch

COPY --from=build /build/mysql-console .

CMD ["/mysql-console"]