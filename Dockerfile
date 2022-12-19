FROM golang:1.19.4-alpine3.17 as builder

# App Name
ARG APP_NAME="search-results-aggregator"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download modules
# Note: prevents layer from being rebuilt unless go.mod and/or go.sum change
COPY go.mod go.sum ./

RUN go mod download

# Copy remaining sources
COPY . .

# Build Service; Extract relevant build output and configurations into separate directory.
RUN go build && \
mkdir -p /out/$APP_NAME && \
cp -r ./$APP_NAME /out/$APP_NAME/ 

FROM alpine:3.17

# App Name
ARG APP_NAME="search-results-aggregator"

WORKDIR /app
# Copy over binaries from alpine builder.
COPY --from=builder /out/$APP_NAME/ ./

EXPOSE 8080
CMD ["./search-results-aggregator"]
