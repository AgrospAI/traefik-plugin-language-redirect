FROM golang:1.25.4-bookworm

ENV GOLANGCI_LINT_VERSION=v2.6.2 \
    YAEGI_VERSION=v0.16.1 \
    CGO_ENABLED=0 \
    GOPATH=/go

ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Install tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl git make ca-certificates unzip bash && \
    rm -rf /var/lib/apt/lists/*

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $GOPATH/bin ${GOLANGCI_LINT_VERSION}

# Install Yaegi
RUN curl -sfL https://raw.githubusercontent.com/traefik/yaegi/master/install.sh | bash -s -- -b $GOPATH/bin ${YAEGI_VERSION}

WORKDIR /workspace/go/src/github.com/your/repo
CMD ["bash"]