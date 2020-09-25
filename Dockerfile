## Step 1: Build the source code
FROM golang:1.15.2-alpine AS builder

WORKDIR $GOPATH/src/jstore
COPY jstore.go jstore.go

# Create non-root appuser
ENV USER=appuser
ENV UID=10001

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Fetch dependencies
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/jstore

# make directory for the data
RUN mkdir -p /store


## Step 2: Make an empty image with the binary
FROM scratch

# Imports from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /store /store

# Copy the static executable
COPY --from=builder /go/bin/jstore /go/bin/jstore

# Set appuser
USER appuser:appuser

# Expose port 8080
EXPOSE 8080

ENTRYPOINT [ "/go/bin/jstore", "-f", "/store", "-u", "http://localhost:8080", "-p", "8080" ]