# ==========================================
# STAGE 1: Compile Go Code into WebAssembly
# ==========================================
FROM golang:1.26-alpine AS builder

# Install basic build tools inside the container
RUN apk add --no-cache git

WORKDIR /src

# Copy module parameters first to leverage Docker caching speeds
COPY go.mod go.sum* ./
RUN go mod download

# Copy the remaining project source files
COPY . .

# Compile the Ebitengine Go source strictly into a WebAssembly binary
RUN env GOOS=js GOARCH=wasm go build -o game.wasm .

# Retrieve the official Javascript runtime bridge relative to this Go version
RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js /src/wasm_exec.js


# ==========================================
# STAGE 2: Package and Serve via Web Server
# ==========================================
FROM nginx:alpine

# Purge default static splash pages from the Nginx container
RUN rm -rf /usr/share/nginx/html/*

# Safely copy compiled components from Stage 1 into the public directory
COPY --from=builder /src/game.wasm /usr/share/nginx/html/
COPY --from=builder /src/wasm_exec.js /usr/share/nginx/html/

# Copy your local custom loader document into the workspace
COPY index.html /usr/share/nginx/html/

EXPOSE 80

