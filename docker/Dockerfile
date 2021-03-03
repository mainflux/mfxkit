FROM golang:1.13-alpine AS builder

WORKDIR /go/src/github.com/mfxkit/mfxkit
COPY . .
RUN apk update \
    && apk add make\
    && apk add git
RUN go get all
WORKDIR /go/src/github.com/mfxkit/mfxkit/cmd/mfxkit
RUN go build main.go -o build/mainflux-mfxkit\
    && mv build/mainflux-mfxkit /exe

FROM scratch
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
