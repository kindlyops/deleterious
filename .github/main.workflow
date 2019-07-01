workflow "Create Release" {
  on = "push"
  resolves = ["goreleaser"]
}

action "is-tag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "goreleaser" {
  uses = "docker://kindlyops/deleterious"
  secrets = [
    "GITHUB_TOKEN",
    "GORELEASER_GITHUB_TOKEN",
  ]
  args = "release"
  needs = ["is-tag"]
}