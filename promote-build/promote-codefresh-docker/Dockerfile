FROM golang:alpine AS build

WORKDIR /app
COPY go/go.mod .
COPY go/go.sum .
RUN go mod download
COPY go/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o promote-codefresh ./src

# ******************************************************

FROM alpine/git:2.36.3

COPY --from=build /app/promote-codefresh /promote-codefresh
ENTRYPOINT ["sh", "-c", "cd ${WORK_DIR:-.}; /promote-codefresh"]
