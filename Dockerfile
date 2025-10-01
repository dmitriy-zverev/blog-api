FROM --platform=linux/arm64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates git curl && rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://go.dev/dl/go1.22.5.linux-arm64.tar.gz | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"
RUN GOMODCACHE=/tmp/go GOCACHE=/tmp/go go install github.com/pressly/goose/v3/cmd/goose@latest

ADD blog-server /usr/bin/blog-server
COPY .env .
COPY sql/schema /migrations

CMD ["sh","-c", "\
   /root/go/bin/goose -dir /migrations postgres $DB_FULL_URL up && \
  exec blog-server" ]
