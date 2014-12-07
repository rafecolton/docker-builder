FROM ubuntu:14.04
MAINTAINER rafael.colton@gmail.com

ENV HOME /root
ENV DEBIAN_FRONTEND noninteractive
ENV GOROOT /usr/local/go
ENV GOPATH /app
ENV GO_TARBALL go1.3.3.linux-amd64.tar.gz
ENV TARBALL_SHA1_SUM 14068fbe349db34b838853a7878621bbd2b24646
ENV LD_LIBRARY_PATH /lib/x86_64-linux-gnu:/usr/local/lib:/usr/lib:/lib

# - Fix some issues with APT packages (See https://github.com/dotcloud/docker/issues/1024)
# - install deps, go
RUN dpkg-divert --local --rename --add /sbin/initctl \
  && ln -sFf /bin/true /sbin/initctl \
  && echo "initscripts hold" | dpkg --set-selections \
  && apt-get update -y \
  && apt-get install -y -qq --no-install-recommends \
    apt-transport-https \
    build-essential \
    curl \
    openssh-client \
    make \
    git-core \
    pkg-config \
    mercurial \
    ca-certificates \
  && update-ca-certificates --fresh \
  && curl -sLO https://storage.googleapis.com/golang/$GO_TARBALL \
  && bash -c "test \"$(openssl sha1 $GO_TARBALL | awk '{print $2}')\" == '$TARBALL_SHA1_SUM'" \
  && tar -C /usr/local -xzf $GO_TARBALL \
  && ln -sv /usr/local/go/bin/* /usr/local/bin \
  && rm -f $GO_TARBALL

WORKDIR /app/src/github.com/rafecolton/docker-builder

ADD Deps /app/src/github.com/rafecolton/docker-builder/Deps
RUN ssh-keyscan github.com > /etc/ssh/ssh_known_hosts \
  && go get github.com/hamfist/deppy \
  && $GOPATH/bin/deppy restore \
  && rm -f $GOPATH/bin/deppy

# set up build dir and add project
ADD . /app/src/github.com/rafecolton/docker-builder

# - make sure we don't have trouble getting deps from GitHub
# - touch Makefile to avoid timestamp error message
# - install
RUN ssh-keyscan github.com > /etc/ssh/ssh_known_hosts \
  && touch Makefile \
  && make build \
  && rm -rf $GOPATH/src \
  && rm -rf $GOPATH/pkg \
  && rm -f $GOPATH/bin/deppy

CMD ["-h"]
ENTRYPOINT ["/app/bin/docker-builder"]
