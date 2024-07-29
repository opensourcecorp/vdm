# vdm: A General-Purpose Versioned-Dependency Manager

`vdm` is an alternative to e.g. Git Submodules for managing arbitrary external
dependencies for the same reasons, in a more sane way. Unlike some other tools
that try to solve this problem, `vdm` is language-agnostic, and can be used for
any purpose that you would need remote development resources.

`vdm` can be used for many different purposes, but most commonly as a way to
track external dependencies that your own code might need, but that you don't
have a language-native way to specify. Some examples might be:

- You have a shared CI repo from which you need to access common shell scripts,
  hosted build tasks, etc.

- You're building & testing a backend application and need to test serving
  frontend code from it, and your team has that frontend code in another
  repository.

- Your team uses protocol buffers and you need to be able to import other loose
  `.proto` files to generate your own code.

`vdm` lets you clearly specify all those remote dependencies & more, and
retrieve them whenever you need them.

## Getting Started

### Installation

`vdm` can be installed from [its GitHub Releases
page](https://github.com/opensourcecorp/vdm/releases). There is a zipped binary
for major platforms & architectures, and those are indicated in the Asset file
name. For example, if you have an M2 macOS laptop, you would download the
`vdm_darwin-arm64.tar.gz` file, and extract it to somewhere on your `$PATH`.

If you have a recent version of the Go toolchain available, you can also install
or run `vdm` using `go`:

```sh
go install github.com/opensourcecorp/vdm@<vX.Y.Z|latest>
# or
go run github.com/opensourcecorp/vdm@<vX.Y.Z|latest> ...
```

### Usage

To get started, you'll need a `vdm` spec file, which is just a YAML (or JSON)
file specifying all your external dependencies along with (usually) their
revisions & where you want them to live on your filesystem:

```yaml
remotes:

  - type:       "git" # the default, and so can be omitted if desired
    remote:     "https://github.com/opensourcecorp/go-common" # can specify as 'git@...' to use SSH instead
    local_path: "./deps/go-common"
    version:    "v0.2.0" # tag example; can also be a branch, commit hash, or the word 'latest'

  - type:       "file" # the 'file' type assumes the version is in the remote field itself somehow, so 'version' can be omitted
    remote:     "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto"
    local_path: "./deps/proto/http/http.proto"
```

You can have as many dependency specifications in that array as you want, and
they can be stored wherever you want. By default, this spec file is called
`vdm.yaml` and lives at the calling location (which is probably your repo's
root), but you can call it whatever you want and point to it using the
`--spec-file` flag to `vdm`.

Once you have a spec file, just run:

```sh
vdm sync
```

and `vdm` will process the spec file, retrieve your dependencies as specified,
and put them where you told them to go. By default, `vdm sync` also removes the
local `.git` directories for each `git` remote, so as to not upset your local
Git tree. If you want to change the version/revision of a remote, just update
your spec file and run `vdm sync` again.

After running `vdm sync` with the above example spec file, your directory tree
would look something like this:

```txt
./vdm.yaml
./deps/
    go-common/
        <stuff in that repo>
    http.proto
```

## Dependencies

`vdm` is distributed as a statically-linked binary per platform that has no
language-specific dependencies. However, note that at the time of this writing,
`vdm` *does* depends on `git` being installed if you specify any `git` remote
types. `vdm` will fail with an informative error if it can't find `git` on your
`$PATH`.

## A note about auth

`vdm` has zero goals to be an authN/authZ manager. If a remote in your spec file
depends on a certain auth setup (an SSH key, something for HTTP basic auth like
a `.netrc` file, an `.npmrc` config file, etc.), that setup is out of `vdm`'s
scope. If required, you will need to ensure proper auth is configured before
running `vdm` commands.

## Future work

- Make the sync mechanism more robust, such that if your spec file changes to
  remove remotes, they'll get cleaned up automatically.

- Add `--keep-git-dir` flag so that `git` remote types don't wipe the `.git`
  directory at clone-time.

- Support more than just `git` and `file` types, and make `file` better
