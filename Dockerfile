FROM golang:1.14-alpine AS build-env
ENV CGO_ENABLED 0
RUN apk add --no-cache git
ADD terracloud /go/src/terracloud

# Install revel framework
RUN go get -u github.com/revel/revel
RUN go get -u github.com/revel/cmd/revel
RUN go get -u github.com/go-playground/validator
RUN go get -u github.com/hashicorp/go-tfe
RUN go get -u github.com/lib/pq
RUN go get -u github.com/iancoleman/strcase
#build revel app
RUN revel build terracloud -m dev terracloud

# Final stage
FROM ubuntu:18.04
RUN apt update && apt upgrade -y && apt install -y apt-transport-https ca-certificates
EXPOSE 9000
WORKDIR /
COPY --from=build-env /go/terracloud /
ENTRYPOINT /run.sh