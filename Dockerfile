FROM golang:1.15-alpine as builder

WORKDIR /go/src/hauk-snitch
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/hauk-snitch .


FROM gcr.io/distroless/base

# Copy our static executable.
COPY --from=builder /go/bin/hauk-snitch /go/bin/hauk-snitch
COPY --from=builder /go/src/hauk-snitch/config.toml /go/src/hauk-snitch/config.toml
WORKDIR /go/src/hauk-snitch
CMD ["/go/bin/hauk-snitch"]