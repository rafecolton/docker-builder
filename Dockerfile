FROM modcloth/docker-builder-base:latest
MAINTAINER devops@modcloth.com

ENV GOPATH /app

## set up build dir and add project
ADD . /app/src/github.com/modcloth/docker-builder
WORKDIR /app/src/github.com/modcloth/docker-builder

# make sure we don't have trouble getting deps from GitHub
RUN ssh-keyscan github.com > /etc/ssh/ssh_known_hosts

# install and verify
RUN touch Makefile
RUN make build

ADD .docker-builder-wrapper /usr/local/bin/docker-builder-wrapper
RUN chmod +x /usr/local/bin/docker-builder-wrapper

EXPOSE 5000
CMD ["-h"]
ENTRYPOINT ["/app/bin/docker-builder"]
