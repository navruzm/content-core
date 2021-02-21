FROM golang:alpine as builder
RUN apk add git make vips-dev fftw-dev build-base --no-cache \
        --repository http://dl-3.alpinelinux.org/alpine/edge/community \
        --repository http://dl-3.alpinelinux.org/alpine/edge/main
# RUN apk add -U --no-cache ca-certificates
RUN apk add --update nodejs yarn
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN yarn install && yarn build
RUN go mod download
RUN GOOS=linux go build -o main ./sample
FROM alpine
RUN apk add git make vips-dev fftw-dev build-base --no-cache \
        --repository http://dl-3.alpinelinux.org/alpine/edge/community \
        --repository http://dl-3.alpinelinux.org/alpine/edge/main
RUN apk add -U --no-cache ca-certificates
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 8080
CMD ["./main"]
