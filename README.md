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

Once you have a spec file, just run:

    vdm sync

and `vdm` will process the spec file, grab your dependencies, put them where
they belong, and check out the right versions. By default, `vdm sync` also
removes the local `.git` directories for each remote, so as to not upset your
local Git tree. If you want to change the version/revision of a remote, just
update your spec file and run `vdm` again.

If for any reason you want all the deps in the spec file to retain their `.git`
directories (such as if you're using `vdm` to initialize a new computer with
actual repos you'd be working in), you can pass the `-keep-git-dir` flag to `vdm
sync`.

## Future work

- Make the sync mechanism more robust, such that if your spec file changes to
  remove remotes, they'll get cleaned up automatically.

- Support more than just Git -- but I really don't care that much about this
  right now.
