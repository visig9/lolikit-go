# lolikit-go

The [Lolinote 2.0] is a simple data specification for personal note-taking.

This tool kit offer some extra conveniences with command-line interface to helping user to manage their [Lolinote 2.0] repository for daily noting.



## Features

- Listing some note by file's modify time or term-frequency score, and optional run / open it.
- Initialize a lolinote repository.
- Serving a lolinote repository as a HTTP file service.
- Print all note's information in JSON format.



## Download

### Pre-build Files

<https://gitlab.com/visig/lolikit-go/tags>

> note: due to developer didn't have some environments, currently only support linux `386`, `amd64` and `arm` platforms, sorry. You can try to build by yourself for local platform. See below.



### Build from Source

Prepare a golang environment, then:

```bash
go get -d -t gitlab.com/visig/lolikit-go/...
cd $(go env GOPATH)/src/gitlab.com/visig/lolikit-go
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

[serve]
address = ":10204"          # default http address
```



## License

MIT



[Lolinote 2.0]: https://gitlab.com/visig/lolinote-spec/blob/master/spec-2.0.md
