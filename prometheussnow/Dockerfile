# multi-stage docker build: compile -> package
#
# compile stage: build the webhook receiver
FROM golang AS builder
WORKDIR /app
COPY . .
RUN go get . && CGO_ENABLED=0 GOOS=linux go build -o snow .

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi-minimal
#MAINTAINER gatblau <onix@gatblau.org>
LABEL author="SNOW"
ARG UNAME=snow
ENV UID=1000
ENV GID=1000
RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream-rpms --enablerepo=ubi-8-baseos-rpms -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g '$GID' -o '$UNAME' && \
    useradd -M -u '$UID' -g '$GID' '$UNAME' && \
    rm -rf /var/cache/yum && \
    microdnf clean all
USER $UNAME
WORKDIR /app
COPY --from=builder /app/snow /app/config.toml ./
CMD ["./snow"]
EXPOSE 8888/tcp
