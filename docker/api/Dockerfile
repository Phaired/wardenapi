ARG GO_IMAGE_VERSION=1.23-alpine3.20

FROM golang:${GO_IMAGE_VERSION}

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/api .

CMD ["api"]