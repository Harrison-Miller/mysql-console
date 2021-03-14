FROM golang:1.16 as build

WORKDIR /build

COPY . .
RUN go get .
RUN CGO_ENABLED=0 go build -o mysql-querier

FROM scratch

COPY --from=build /build/mysql-querier .

CMD ["/mysql-querier"]