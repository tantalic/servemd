# `servemd`

Servemd makes it simple and fast to serve one or more [markdown][markdown] files without the build process required by static site generators. Servemd can be used for quick local viewing or to serve content to thousands of visitors.

## Installation

Servemd is distributed as either:

- An executable with no dependencies
- A docker image

### Executable

Servemd is a single executable with no dependencies. Installation is as simple as downloading the binary for your platform from the [release page][release]. Releases are made available for:

- Mac OS X (x86, x64)
- Linux (x86, x64, ARMv5, ARMv6, ARMv7, ARMv8)
- FreeBSD (x86, x64, ARMv5, ARMv6, ARMv7)
- NetBSD (x86, x64, ARMv5, ARMv6, ARMv7)
- DragonFly BSD (x64)
- Windows (x86, x64)

OS X users can install `servemd` using the [Homebrew package manager][homebrew]:

```shell
brew tap tantalic/tap
brew install servemd
```

### Docker

The [tantalic/servemd][dockerhub] image can be pulled from Docker Hub:

```shell
docker pull tantalic/servemd
```

## Configuration

Usage: `servemd [OPTIONS] [DIR]`

### Arguments

| Environment Variable | Argument |           Description            | Default Value |
|----------------------|----------|----------------------------------|---------------|
| `DOCUMENT_ROOT`      | `[DIR]`  | Directory to serve content from. | `.`           |

### Flags

| Environment Variable |          Flag          |                                          Description                                          | Default Value |
|----------------------|------------------------|-----------------------------------------------------------------------------------------------|---------------|
| `HOST`               | `-a, --host`           | Host/IP address to listen on                                                                  | All addresses |
| `PORT`               | `-p, --port`           | TCP port to listen on                                                                         | `3000`        |
| `BASIC_AUTH`         | `-u, --auth`           | Username and password for HTTP basic authentication. In the form of `user1:pass1,user2:pass2` | None          |
| `DOCUMENT_EXTENSION` | `-e, --extension`      | Extension used for markdown files                                                             | `.md`         |
| `DIRECTORY_INDEX`    | `-i, --index`          | Filename (without extension) to use for directory indexes                                     | `index`       |
| `MARKDOWN_THEME`     | `-m, --markdown-theme` | Theme to use for styling markdown. Values: *clean*, *github*, *developer*                     | `clean`       |
| `CODE_THEME`         | `-c, --code-theme`     | Syntax highlighting theme (powered by [highlight.js][highlightjs])                            | None          |


## Deploying with Docker 

The [tantalic/servemd][dockerhub] image can be used as a base for deployment. To create an image simply add your content to `/app/content` and configure via environment variables in your `Dockerfile`:

```Dockerfile
FROM tantalic/servemd:latest
MAINTAINER Your Name <email@example.com>

ENV MARKDOWN_THEME developer
ENV CODE_THEME solarized-dark

ADD content /app/content
```


[markdown]: https://daringfireball.net/projects/markdown/syntax
[release]: https://github.com/tantalic/servemd/releases/latest
[homebrew]: http://brew.sh
[highlightjs]: http://highlightjs.org
[dockerhub]: https://hub.docker.com/r/tantalic/servemd/


