FROM golang:1.20 as devimage

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.51.2
WORKDIR /go/src/github.com/LeoCBS/httpmiddleware
COPY go.mod /go/src/github.com/LeoCBS/httpmiddleware
RUN go mod download
RUN go mod tidy
COPY . /go/src/github.com/LeoCBS/httpmiddleware
