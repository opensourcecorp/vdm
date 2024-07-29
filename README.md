# vdm: A General-Purpose Versioned-Dependency Manager

`vdm` is an alternative to e.g. git submodules for managing arbitrary external
dependencies for the same reasons, in a more sane way. Unlike some other tools
that try to solve this problem, `vdm` is language-agnostic, and can be used for
any purpose that you would need remote development resources.

To get started, you'll need a `vdm` spec file, which is just a YAML (or JSON)
file specifying all your external dependencies along with (usually) their
revisions & where you want them to live on your filesystem:

```yaml
remotes:

  - type:       "git" # the default, and so can be omitted if desired
    remote:     "https://github.com/opensourcecorp/go-common"
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

and `vdm` will process the spec file, retrieve your dependencies, put them where
they belong, and check out the right versions. By default, `vdm sync` also
removes the local `.git` directories for each `git` remote, so as to not upset
your local Git tree. If you want to change the version/revision of a remote,
just update your spec file and run `vdm sync` again.

After running `vdm sync` with the above example spec file, your directory tree
would look something like this:

```txt
./vdm.yaml
./deps/
    go-common/
        <stuff in that repo>
    http.proto
```

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

- Support more than just Git
