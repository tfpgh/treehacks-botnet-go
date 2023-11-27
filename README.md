# treehacks-botnet-go

Recreating [TreeHacks/botnet-hackpack](https://github.com/TreeHacks/botnet-hackpack) in [go](https://go.dev/) as my first project in the language.

The original project just includes a binary for the master. I'm rewriting both the bot and master here.

## pre-commit

pre-commit hooks should be installed before commiting.
Instructions for installing pre-commit itself are [here](https://pre-commit.com/#install). To install hooks run `pre-commit install`.

Hook dependencies:

- [go](https://go.dev/) (obviously)
- [go-critic](https://github.com/go-critic/go-critic)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

Hooks can be run manually with `pre-commit run --all-files`
