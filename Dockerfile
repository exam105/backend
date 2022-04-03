FROM golang:1.16.5 AS builderStep

LABEL author="Muhammad Tariq"

ARG Mongo_User=testuser
ENV ENV_MONGO_USER=$Mongo_User

RUN echo "Environment Variable:=> $ENV_MONGO_USER"
RUN echo "Mongo User ARG Variable:=> $Mongo_User"

ENV APP_HOME /go/src/github.com/exam105-UPD/backend

# Make Build dir
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

# This change is made for fixing S3 UploadImage to Dev Server
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

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

# This change is made for fixing S3 UploadImage to Dev Server
COPY --from=builderStep /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builderStep /go/src/github.com/exam105-UPD/backend .
EXPOSE 9090
ENTRYPOINT [ "./exam105" ]
