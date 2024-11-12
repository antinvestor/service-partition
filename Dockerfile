FROM golang:1.23 as builder

WORKDIR /

COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o partition_binary .

FROM gcr.io/distroless/static:nonroot

USER nonroot:nonroot
EXPOSE 80
EXPOSE 50051

WORKDIR /

COPY --from=builder /partition_binary /partition
COPY --from=builder /migrations /migrations

# Run the service command by default when the container starts.
ENTRYPOINT ["/partition"]
