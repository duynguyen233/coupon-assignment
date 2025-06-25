# syntax=docker/dockerfile:1
# Stage 1 - get env
FROM alpine:3.18 as certificates
ARG GITLAB_TOKEN
RUN echo I received "${GITLAB_TOKEN}" from GitLab!
ARG DEPLOY_TAG  
RUN echo I received "${DEPLOY_TAG}" from GitLab!
ARG CACHEBUST
RUN echo CACHEBUST "$CACHEBUST"
RUN apk --no-cache add curl
RUN curl --header "PRIVATE-TOKEN: ${GITLAB_TOKEN}" "https://gitlab.rockship.co/api/v4/projects/16/repository/files/backend%2Fcarta%2F${DEPLOY_TAG}%2Fconfig.yml/raw?ref=master" > config.yml

FROM golang:latest as builder

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .
COPY --from=certificates config.yml config/config.yml

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
FROM scratch
ENV ENV PROD
COPY --from=builder docker-gs-ping docker-gs-ping
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certificates config.yml config/config.yml
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]

