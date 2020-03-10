FROM golang:1.7.3 AS builder

RUN mkdir /build
COPY ./ /build
WORKDIR /build
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

WORKDIR /
COPY --from=builder /build/main .
COPY --from=builder /build/.env .

CMD ["/main"]
