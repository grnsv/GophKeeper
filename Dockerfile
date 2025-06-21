FROM golang:1.24-alpine AS builder

ARG VERSION
ARG DATE
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goph-keeper-server \
    -ldflags "-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${DATE}'" \
    ./cmd/server


FROM scratch

COPY --from=builder /go/src/goph-keeper-server /goph-keeper-server
COPY ./migrations /migrations
CMD ["/goph-keeper-server"]
