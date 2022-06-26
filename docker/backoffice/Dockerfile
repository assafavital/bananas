# syntax=docker/dockerfile:1.3
FROM golang:bullseye as builder

WORKDIR /code/
COPY ./go.mod ./go.sum /code/
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY Makefile /code/
COPY pkg /code/pkg
COPY cmd/backoffice /code/cmd/backoffice

ARG GIT_TAG=""
ARG GIT_COMMIT=""
ARG RC_TARGET=""
# hadolint ignore=DL3509
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build/linux/backoffice NOPROTOGEN=1 GOOS=linux GOARCH=amd64 GIT_TAG=$GIT_TAG GIT_COMMIT=$GIT_COMMIT RC_TARGET=$RC_TARGET

FROM alpine:3.16.0
# Add AWS root certificates to trusted roots for RDS access
RUN wget -nv -O - https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem >> /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /code/build/linux/backoffice /server
RUN mkdir -p /var/local/certs
ENTRYPOINT ["/server"]
CMD [ "run" ]
