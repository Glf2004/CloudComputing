FROM golang:1.24.2-bookworm AS builder
WORKDIR /gomall
COPY . .
# set GOPROXY to use the proxy server in China
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN chmod +x /gomall/app/email/build.sh && cd /gomall/app/email && ./build.sh

FROM debian:stable-slim AS runtime
WORKDIR /app
COPY --from=builder /gomall/app/email/output /app
CMD ["bash", "/app/bootstrap.sh"]