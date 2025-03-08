# FROM golang:1.23 AS builder
# #FROM harbor.ks.x/eclipse-che/golang:1.22 as builder
# # ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# ENV CGO_ENABLED=0
# ENV GOPROXY="https://goproxy.cn"
# COPY . /go/src/
# WORKDIR /go/src/
# RUN unset GOPATH && pwd && ls -l && \
#     go build -mod=vendor -v

# FROM harbor.ks.x/devops/alpine@sha256:93d5a28ff72d288d69b5997b8ba47396d2cbb62a72b5d87cd3351094b5d578a0
# FROM alpine:3.21
FROM ubuntu:20.04
# FROM harbor.ks.x/eclipse-che/ubuntu:20.04 
# COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
# ENV ZONEINFO /opt/zoneinfo.zip
# ADD config.yaml /config.yaml
# COPY --from=builder /go/src/showme /showme
COPY showme /showme
# COPY apis /apis
EXPOSE 8999
# CMD ["sh", "-c", "/aliveagent -server"]
ENTRYPOINT [ "/showme" ]
