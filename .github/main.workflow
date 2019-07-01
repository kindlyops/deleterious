workflow "Create Release" {
  on = "push"
  resolves = ["goreleaser"]
}

action "is-tag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "goreleaser" {
  uses = "kindlyops/deleterious@master"
  secrets = [
    "GITHUB_TOKEN",
  ]
  args = "release"
  needs = ["is-tag"]
}