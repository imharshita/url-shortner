# Build backend exec file.
FROM golang:1.22-alpine AS builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o main .

# Make workspace with above generated files.
FROM gcr.io/distroless/static:nonroot

COPY  . /app

COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 3030

CMD ["./main"]