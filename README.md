# lolikit-go

[![Build Status](https://travis-ci.org/visig9/lolikit-go.svg?branch=master)](https://travis-ci.org/visig9/lolikit-go)

The [Lolinote 2.0] is a simple data specification for personal note-taking.

This tool kit offer some extra conveniences with command-line interface to helping user to manage their [Lolinote 2.0] repository for daily noting.

> This program still in heavy developing and the interface may be changed in future.



## Features

- Initialize a lolinote repository.
- Serving a lolinote repository as a HTTP file service.
- Listing some note by file's modify time or term-frequency score, and optional run / open it.
- Print all metadata of notes in JSON format for further processing.



## Download

### Pre-build Files

<https://github.com/visig9/lolikit-go/releases>

> note: due to developer didn't have some environments, currently only support linux `386`, `amd64` and `arm` platforms, sorry. You can try to build by yourself for local platform. See below.



### Build from Source

Prepare a golang environment, then:

```bash
go get -d -t github.com/visig9/lolikit-go/...
cd $(go env GOPATH)/src/github.com/visig9/lolikit-go
./maintain.sh install
```



## Configuration

Lolikit using 2 configuration files at the same time. It's *User's* and *Repository's* configuration file. Those paths are...

- User's configuration file:
    - if `$XDG_CONFIG_HOME` exists, using `$XDG_CONFIG_HOME/lolikit/config.toml`.
    - else `$HOME/.config/lolikit/config.toml`.
- Repository's configuration file:
    - `<repo-dir>/.lolinote/lolikit/config.toml`

The *Repository's* configuration will shadow the *User's* one in each options level.

Here is a configuration example.

```toml
# The default repository. Only useful in User's configuration file.
default-repo = "/path/to/my/lolinote/repo"

# "text-types" setting only using for searching functions.
# If a note's content-type is one of the text-types, the relevance
# calculation will try to anaylze both the content and filename. Else, only
# the filename.
# The content in files will always be treated as utf-8 encoding.
text-types = [
    "txt",
    "md",
]

[list]
page-size = 10              # default page size of a list
runner = "xdg-open"         # default runner for any content type
dir-runner = "nautilus"     # default directory runner

[list.runners]              # default runner for particular content-type
md  = "wc -l"
txt = "gedit"
jpg = "firefox"

[new]
buffer = "to-be-classfied"  # the buffer area for quick added note
buffer-size = 10            # the maximum entry should take time to moving.

[serve]
address = ":10204"          # default http address
```



## License

MIT



[Lolinote 2.0]: https://github.com/visig9/lolinote-spec
