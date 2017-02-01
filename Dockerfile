FROM scratch
MAINTAINER Kevin Stock <kevinstock@tantalic.com>

WORKDIR /content
VOLUME /content

ADD ./build/servemd-linux_amd64 /servemd
EXPOSE 3000
CMD ["/servemd"]