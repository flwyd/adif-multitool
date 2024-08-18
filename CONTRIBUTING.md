# How to Contribute

We'd love to accept your patches and contributions to this project. There are
just a few small guidelines you need to follow.

## Follow the code of conduct

Open Source works best when everyone feels comfortable and empowered to
participate. Please read and follow the [code of conduct](CODE_OF_CONDUCT.md)
and don't be a jerk.

## Discuss new functionality

This project takes an opinionated approach, following the Unix pipeline
philosophy. There are many features that adif-multitool could add, but they
should be done in harmony with the tool's philosophy. For new features or
changes in behavior, please open an issue first to discuss the semantics of the
feature. For straightforward bug fixes, a pull request by itself is sufficient.

## Add tests

It helps when bug fixes include a test case that fails without the fix.
New functionality should be covered by automated tests to avoid regressions.
Commands and many API functions have standard Go unit tests.  There are also
some tests that exercise the command line, including flag syntax parsing and
checking stderr, in the [txtar](https://pkg.go.dev/golang.org/x/tools/txtar)
files [in the `adifmt` package](./adifmt/testdata).  These command language for
these tests is described in the
[testscript package](https://pkg.go.dev/github.com/rogpeppe/go-internal@v1.12.0/testscript).
