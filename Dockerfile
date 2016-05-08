FROM gliderlabs/alpine
MAINTAINER Kevin Stock <kevinstock@tantalic.com>
RUN apk-install ca-certificates 

RUN mkdir /app
RUN mkdir /app/content

ADD ./servemd /app
RUN chmod a+x /app/servemd

ENV PORT 8000
WORKDIR /app/content
ENTRYPOINT /app/servemd
EXPOSE 8000
