## What this changes

<!-- The behaviour difference, in a sentence or two. -->

## Why

<!-- What problem this solves. If it fixes a bug, describe how the bug manifests. -->

## Where does this belong?

This repository is a thin Go binding. Most behaviour lives in
[phenixrizen/zen](https://github.com/phenixrizen/zen), and `deps/` holds static libraries built
from it.

- [ ] Go surface only — no engine change needed
- [ ] Requires a matching change in `phenixrizen/zen` (link it here)
- [ ] Requires rebuilt `deps/` libraries

## Testing

- [ ] `make fmt_check`
- [ ] `make test`
- [ ] For a bug fix: the new test fails without the fix
