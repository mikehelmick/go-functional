# go-functional - contribution guidelines

Contributions to this project are welcomed. We request that you
read through the guidelines before getting started.

## Contributor License Agreement

Contributions to this project must be accompanied by a Contributor License
Agreement (CLA). You (or your employer) retain the copyright to your
contribution; this simply gives us permission to use and redistribute your
contributions as part of the project.

To sign the CLA, open a signed PR adding yourself as a contributor
in the .clabot file at the root of this repository. Before doing so, please
review the cla.md file.

## Community guidelines

TBD.

## Contribution

### Code reviews

All submissions will be reviewed before merging. Submissions are reviewed using
[GitHub pull requests](https://help.github.com/articles/about-pull-requests/).

Please note that the `federation*` packages are reference-only, and we do not
actively support them.

## Source and build

### Source code layout

This project depends in Go 1.20 or newer.

### Running tests

Run the tests with:

```text
$ go test ./...
```

### Presubmit checks
