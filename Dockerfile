FROM rafecolton/docker-builder-base:latest
MAINTAINER rafael.colton@gmail.com

ENV GOPATH /app

## set up build dir and add project
ADD . /app/src/github.com/rafecolton/docker-builder
WORKDIR /app/src/github.com/rafecolton/docker-builder

# make sure we don't have trouble getting deps from GitHub
RUN ssh-keyscan github.com > /etc/ssh/ssh_known_hosts

# install and verify
RUN touch Makefile
RUN make build

EXPOSE 5000
CMD ["-h"]
ENTRYPOINT ["/app/bin/docker-builder"]
