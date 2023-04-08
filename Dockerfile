# Building the binary of the App
FROM golang:1.18 as build

# Define build env
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
#ENV GOARCH=amd64

# `beer-api` should be replaced with your project name
WORKDIR /api

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN go build -o app

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest

RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates tzdata ffmpeg && \
  rm -rf /var/cache/*

# Create the `public` dir and copy all the assets into it
COPY ./configs ./configs
COPY ./templates ./templates

COPY --from=build /api/app .

EXPOSE 8000

CMD ["./app","-environment", "dev"]