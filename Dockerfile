FROM golang:1.13.8 AS builderStep

LABEL author="Muhammad Tariq"

# Make Build dir
WORKDIR /src

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
WORKDIR /app
COPY --from=builderStep /src .
EXPOSE 9090
ENTRYPOINT [ "./exam105" ]
