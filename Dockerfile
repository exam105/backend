FROM golang:1.13.8 AS builderStep

LABEL author="Muhammad Tariq"

ENV APP_HOME /go/src/github.com/exam105-UPD/backend

# Make Build dir
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

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
WORKDIR /go/src/github.com/exam105-UPD/backend
COPY --from=builderStep /go/src/github.com/exam105-UPD/backend .
EXPOSE 9090
ENTRYPOINT [ "./exam105" ]
