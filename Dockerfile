FROM golang:1.12 AS builder

RUN mkdir /build
COPY ./ /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

WORKDIR /
COPY --from=builder /build/main .

CMD ["/main"]