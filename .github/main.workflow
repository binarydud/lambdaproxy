workflow "Test" {
  on = "pull_request"
  resolves = ["lint", "test"]
}

action "lint" {
  uses = "apex/actions/go@master"
  args = "make lint"
}

action "test" {
  uses = "apex/actions/go@master"
  args = "make test"
}
