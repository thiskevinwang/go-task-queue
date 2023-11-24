# Use an official Golang runtime as a parent image
FROM golang:1.21.3-bookworm as builder
WORKDIR /app
COPY /go.work /app
COPY /go.work.sum /app
RUN go work edit -dropuse=./apps/server

COPY /apps/runner /app/apps/runner
COPY /packages/shared /app/packages/shared
# RUN --mount=type=cache,target=/Users/me/Library/Caches go mod download

# output binary named "main"
RUN go build -o main apps/runner/runner.go


FROM golang:1.21.3-bookworm 
WORKDIR /app
COPY --from=builder /app/main /app/main


# Run the Go app when the container starts
CMD ["./main"]
