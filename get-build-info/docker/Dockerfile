FROM golang:alpine AS build

WORKDIR /app
COPY go/go.mod .
COPY go/go.sum .
RUN go mod download
COPY go/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o metadata ./src

# ******************************************************

FROM alpine/git:2.36.3

COPY --from=build /app/metadata /metadata
ENTRYPOINT "ls -la && /metadata meta export"
