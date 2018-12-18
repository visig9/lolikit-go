# lolikit-go

The [Lolinote 2.0] is a simple data specification for personal note-taking.

Lolikit offer some extra conveniences to help users to manage their data for daily noting.

> Hint: Lolinote wasn't lean on this program, just for conveniences.



## Configuration

Lolikit using 2 configuration files at the same time. It called *User* and *Repository* configuration file. Those paths are...

- *User* configuration file:
    - if `$XDG_CONFIG_HOME` exists, using `$XDG_CONFIG_HOME/lolikit/config.toml`.
    - else `$HOME/.config/lolikit/config.toml`.
- *Repository* configuration file:
    - `<repo-dir>/.lolinote/lolikit/config.toml`

The *Repository* configuration will shadow the *User's* configuration in each options level.

Here is a configuration example.

```toml
# The default repository, only useful in User configuration file.
default-repo = "/path/to/my/lolinote/repo"

# "text-types" setting only using for search relative functions.
# If a note's content-type is one of the text-types, the relevance
# calculation will try to anaylze both the content and filename. Else, only
# the filename.
# The content in files will be treat as utf-8 encoding.
text-types = [
    "txt",
    "md",
]

[list]
page-size = 10              # default page size of a list
runner = "xdg-open"         # default runner for any content type
dir-runner = "nautilus"     # default directory runner

[list.runners]         # default runner for particular content-type
md  = "wc -l"
txt = "gedit"
jpg = "firefox"
```



## Download

### Pre-build Files

<https://gitlab.com/visig/lolikit-go/tags>

> note: due to developer didn't have some environments, currently only support linux `386`, `amd64` and `arm` platforms, sorry.



### Build from Source

Prepare a golang environment, then:

```bash
go get -d -t gitlab.com/visig/lolikit-go/...
cd $(go env GOPATH)/src/gitlab.com/visig/lolikit-go
./maintain.sh install
```


## License

MIT



[Lolinote 2.0]: https://gitlab.com/visig/lolinote-spec/
