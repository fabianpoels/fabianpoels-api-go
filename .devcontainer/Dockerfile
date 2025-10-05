FROM golang:1.25.1-bookworm

# Set working directory
WORKDIR /app

# Install only essential system tools
RUN apt-get update && apt-get install -y \
    redis-tools \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid 1000 -m vscode \
    && chown -R vscode:vscode /app

# Switch to the vscode user for the remaining operations
USER vscode

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/air-verse/air@latest

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

# Default command when container starts
CMD ["bash"]