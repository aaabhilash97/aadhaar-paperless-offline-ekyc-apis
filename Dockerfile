FROM golang:1.17.2-bullseye as build-env

# All these steps will be cached
RUN mkdir /app-server
WORKDIR /app-server
# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .
ARG CI_COMMIT_SHA
ARG CI_COMMIT_REF_NAME
# RUN MKDIR -p dist
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.gitCommit=$CI_COMMIT_SHA -X main.gitRef=$CI_COMMIT_REF_NAME" -a -installsuffix cgo -o ./app ./cmd/server/main.go
RUN ./app -version

# <- Second step to build minimal image
FROM gcr.io/distroless/static-debian11

WORKDIR /app-server

COPY --from=build-env /app-server/app /app-server

ENTRYPOINT ["/app-server/app"]
EXPOSE 3333
EXPOSE 4444