FROM golang:1.16-buster as builder

WORKDIR /

COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o partition_binary .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /partition_binary /partition
COPY --from=builder /migrations /migrations

WORKDIR /

# Run the service command by default when the container starts.
ENTRYPOINT ["/partition"]

# Document the port that the service listens on by default.
EXPOSE 7003