# `multidockerfile`

`multidockerfile` is a command-line tool that allows you to split multi-stage Dockerfiles into multiple files.


## Usage

### `multidockerfile join <inputs> ...`

Join multiple Dockerfiles into a single multi-stage Dockerfile.

#### Arguments

| Name | Description |
|-|-|
| `<inputs> ...` | Paths to the Dockerfiles to be joined. |

#### Options

| Short | Long | Default | Description |
|-|-|-|-|
| `-o` | `--output` | `-` | Where to write the multi-stage Dockerfile (`-` for stdout). |

#### Example

```dockerfile
# dockerfiles/one.dockerfile

FROM alpine AS one
```

```dockerfile
# dockerfiles/two.dockerfile

FROM alpine AS two
```

```console
$ multidockerfile join dockerfiles/*.dockerfile
FROM alpine AS one
FROM alpine AS two
```

### `multidockerfile version`

Show the `multidockerfile` version information.

#### Example

```console
$ multidockerfile version
{
  "Version": "0.1.0-dev",
  "GitCommit": "20586c3eb00aad3dde1ca63eb47dcb14ae6372d5",
  "Built": "2021-02-26T21:22:04Z",
  "GoVersion": "go1.16",
  "OS": "darwin",
  "Arch": "amd64"
}
```
