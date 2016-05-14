FROM gliderlabs/alpine
MAINTAINER Kevin Stock <kevinstock@tantalic.com>

RUN mkdir /app
RUN mkdir /app/content

ADD ./servemd /app
RUN chmod a+x /app/servemd

WORKDIR /app/content
ENTRYPOINT /app/servemd
EXPOSE 3000
