FROM golang:1.23.5 as builder

WORKDIR /dlza-manager-storage-handler

ARG SSH_PUBLIC_KEY=$SSH_PUBLIC_KEY
ARG SSH_PRIVATE_KEY=$SSH_PRIVATE_KEY

# ARG GITLAB_USER=gitlab-ci-token
# ARG GITLAB_PASS=$CI_JOB_TOKEN
# ARG SSH_PRIVATE_KEY
# ARG SSH_PUBLIC_KEY

ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.switch.ch/ub-unibas/*
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

# RUN cat go.mod
RUN apt-get update && \
    apt-get install -y \
        git \
        openssh-client \
        ca-certificates \
        protobuf-compiler 
# RUN apk add --no-cache ca-certificates git openssh-client 
# RUN 'which ssh-agent || ( apt-get update -y && apt-get install openssh-client git -y )'
RUN eval $(ssh-agent -s)
RUN mkdir -p ~/.ssh
RUN chmod 700 ~/.ssh
# #for CI/CD build
RUN echo "$SSH_PRIVATE_KEY" | base64 -d >> ~/.ssh/id_rsa
# #for local build
# RUN echo "$SSH_PRIVATE_KEY" >> ~/.ssh/id_rsa
RUN echo "$SSH_PUBLIC_KEY" | tr -d '\r'   >> ~/.ssh/authorized_keys
# # set chmod 600 else bas permission it fails
RUN chmod 600 ~/.ssh/id_rsa
RUN chmod 644 ~/.ssh/authorized_keys
RUN ssh-keyscan gitlab.switch.ch >> ~/.ssh/known_hosts
RUN chmod 644 ~/.ssh/known_hosts
# RUN git config --global url."ssh://git@gitlab.switch.ch/".insteadOf "https://gitlab.switch.ch/"
RUN git config --global --add url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.switch.ch".insteadOf "https://gitlab.switch.ch"
# RUN ssh -A -v -l git gitlab.switch.ch

# with DOCKER_BUILDKIT=1 for ssh
# RUN --mount=type=ssh go mod download
RUN go mod download
# RUN git clone https://${GITLAB_USER}:${GITLAB_PASS}@gitlab.switch.ch/ub-unibas/dlza/microservices/pbtypes /pbtypes
# RUN go get google.golang.org/protobuf/protoc-gen-go 
# RUN go get google.golang.org/protobuf
# RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
# RUN bash ./build.sh
RUN go build 

FROM scratch
WORKDIR /
COPY --from=builder /dlza-manager-storage-handler /
EXPOSE 8080

ENTRYPOINT ["/dlza-manager-storage-handler"]