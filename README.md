# treehacks-botnet-go

Following [TreeHacks/botnet-hackpack](https://github.com/TreeHacks/botnet-hackpack) as my first project in [go](https://go.dev/).

## pre-commit

pre-commit hooks should be installed before commiting.
Instructions for installing pre-commit itself are [here](https://pre-commit.com/#install). To install hooks run `pre-commit install`.

Hook dependencies:

- [go](https://go.dev/) (obviously)
- [go-critic](https://github.com/go-critic/go-critic)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

Hooks can be run manually with `pre-commit run --all-files`
