load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "rules_rust",
    sha256 = "872b04538ca20dad94791c348623f079ba93daf274c1d57ae6bfe0930ec77f0d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_rust/releases/download/0.6.0/rules_rust-v0.6.0.tar.gz",
        "https://github.com/bazelbuild/rules_rust/releases/download/0.6.0/rules_rust-v0.6.0.tar.gz",
    ],
)

load("@rules_rust//rust:repositories.bzl", "rules_rust_dependencies", "rust_register_toolchains")

rules_rust_dependencies()

rust_register_toolchains(version = "1.61.0", edition="2021")

load("@rules_rust//proto:repositories.bzl", "rust_proto_repositories")

rust_proto_repositories()

load("@rules_rust//proto:transitive_repositories.bzl", "rust_proto_transitive_repositories")

rust_proto_transitive_repositories()
