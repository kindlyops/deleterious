# Bazel workspace created by @bazel/create 0.37.1

# Declares that this directory is the root of a Bazel workspace.
# See https://docs.bazel.build/versions/master/build-ref.html#workspace
workspace(
    # How this workspace would be referenced with absolute labels from another workspace
    name = "com_kindlyops_deleterious",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Keep the version of rules_go in sync with go.mod as well, so that dependabot
# will tell us about new releases of rules_go
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "86d3dc8f59d253524f933aaf2f3c05896cb0b605fc35b460c0b4b039996124c6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.22.8")

http_archive(
    name = "bazel_gazelle",
    sha256 = "e467b801046b6598c657309b45d2426dc03513777bd1092af2c62eebf990aca5",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.45.0/bazel-gazelle-v0.45.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.45.0/bazel-gazelle-v0.45.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

buildtools_version = "4.2.3"

http_archive(
    name = "io_bazel_buildtools",
    sha256 = "5ec71602e9b458b01717fab1d37492154c1c12ea83f881c745dbd88e9b2098d8",
    strip_prefix = "buildtools-{0}".format(buildtools_version),
    urls = ["https://github.com/bazelbuild/buildtools/archive/{0}.tar.gz".format(buildtools_version)],
)

http_archive(
    name = "rules_pkg",
    sha256 = "d250924a2ecc5176808fc4c25d5cf5e9e79e6346d79d5ab1c493e289e722d1d0",
    url = "https://github.com/bazelbuild/rules_pkg/releases/download/0.10.1/rules_pkg-0.10.1.tar.gz",
)

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()
