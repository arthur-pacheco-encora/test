# Use the official golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM us-east1-docker.pkg.dev/anchorage-ci/docker-build/golang:1.22-bullseye@sha256:0c099cf97c6162ceab8c9494c4ce8c2a19a73b852f334ec1876a42d05eb5f90b AS builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binaries
RUN CGO_ENABLED=0 go build -v -o bin/ ./...

FROM us-east1-docker.pkg.dev/anchorage-ci/docker-build/google/cloud-sdk:alpine@sha256:1ef162f30c0da2b61277ec135335ec4c23d94a6a87c443d6b9f7b4458ddde532

# Environment variables for the project
ENV PROJECT_ID="development-204920"
ENV SHEETS_API_BASE_URL="https://sheets.googleapis.com/v4/spreadsheets"


# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/server /server

# Run the web service on container startup.
CMD ["/server"]
