# `multidockerfile`

`multidockerfile` is a command-line tool that allows you to split multi-stage Dockerfiles into multiple files.

## Background

[Multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/) are a powerful way to optimize Dockerfiles, especially when combined with [`docker buildx bake`](https://github.com/docker/buildx/blob/master/docs/reference/buildx_bake.md) to build multiple images in parallel.
However, Docker requires that all stages are defined in a single Dockerfile, which can become difficult to navigate as it grows.

With `multidockerfile`, you can split the Dockerfile up, and recombine it with `multidockerfile join` before building the images.

`multidockerfile` parses the individual Dockerfiles looking for the `FROM` and `COPY --from` instructions that create dependencies between stages.
The combined Dockerfile is sorted so that stages with dependencies appear after the stages on which they depend.

## Installation

### From binary releases

Pre-built binaries are available for each [release](https://github.com/haines/multidockerfile/releases).
You can download the correct version for your operating system, make it executable with `chmod +x`, and either execute it directly or put it on your path.

SHA-256 checksums and GPG signatures are available to verify integrity.
My GPG public key can be obtained from

<details>
  <summary>GitHub (<a href="https://github.com/haines">@haines</a>)</summary>

  ```console
  $ curl https://github.com/haines.gpg | gpg --import
  ```
</details>

<details>
  <summary>Keybase (<a href="https://keybase.io/haines">haines</a>)</summary>

  ```console
  $ curl https://keybase.io/haines/pgp_keys.asc | gpg --import
  ```
</details>

<details>
  <summary>keys.openpgp.net (<a href="https://keys.openpgp.org/search?q=andrew%40haines.org.nz">andrew@haines.org.nz</a>)</summary>

  ```console
  $ gpg --keyserver keys.openpgp.org --recv-keys 6E225DD62262D98AAC77F9CDB16A6F178227A23E
  ```
</details>

### With Docker

A Docker image is available for each release at [ghcr.io/haines/multidockerfile](https://ghcr.io/haines/multidockerfile).

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
