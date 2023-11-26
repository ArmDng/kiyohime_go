FROM golang:latest AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /kiyohime

FROM build-stage AS run-stage

RUN go test -v ./...

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /kiyohime /kiyohime

RUN apk --no-cache add tzdata


ENTRYPOINT [ "./kiyohime" ]