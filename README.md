# vdm: Versioned Dependency Manager

`vdm` is an alternative to git submodules for managing external dependencies for
the same reasons, in a more sane way.

To get started, you'll need a `vdm` spec file, which is just a JSON array of all
your external dependencies along with their revisions & where you want them to
live in your repo:

    [
      {
        "remote":     "https://github.com/opensourcecorp/go-common",
        "version":    "v0.2.0", // tag; can also be a branch, short or long commit hash, or the word 'latest'
        "local_path": "./deps/go-common"
      }
    ]

You can have as many dependency specifications in that array as you want. By
default, this spec file is called `vdm.json` and lives at the calling location
(which is probably your repo's root), but you can call it whatever you want and
point to it using the `-spec-file` flag to `vdm`.

Then, just run:

    vdm

and `vdm` will process the spec file, grab your dependencies, put them where
they belong, and check out the right versions. `vdm` also removes the local
`.git` directories for each remote, so as to not upset your local Git tree. If
you want to change the version/revision of a remote, just update your spec file
and run `vdm` again.

## Future work

- Make keeping the `.git` directory optional, so people can use it as an "all my
  repos on a fresh machine" manager. I know I would want this lol.

- Have some kind of sync mechanism, such that if your spec file changes to
  remove remotes, they'll get cleaned up automatically.

- Support more than just Git -- but I really don't care that much about this
  right now.
