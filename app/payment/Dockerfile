FROM golang:1.24.2-bookworm AS builder
WORKDIR /gomall
COPY . .
# set GOPROXY to use the proxy server in China
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN chmod +x /gomall/app/payment/build.sh && cd /gomall/app/payment && ./build.sh

FROM debian:stable-slim AS runtime
WORKDIR /app
COPY --from=builder /gomall/app/payment/output /app
CMD ["bash", "/app/bootstrap.sh"]