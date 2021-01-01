FROM golang:1.13.8 AS builderStep

# Make Build dir
RUN mkdir /go/src/exam105/backend
WORKDIR /go/src/exam105/backend

# Copy golang dependency manifests
COPY go.mod .
COPY go.sum .

# Cache the downloaded dependency in a layer.
RUN go mod download

# add the source code
COPY . .

# Build
RUN go get && CGO_ENABLED=0 GOOS=linux go build -o exam105 .

FROM scratch AS app
COPY --from=builderStep /go/src/exam105/backend /opt/exam105
WORKDIR /opt/exam105/
EXPOSE 9090
ENTRYPOINT [ "./exam105" ]
