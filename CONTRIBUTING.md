# Contributing

Thank you for considering contributing to the Pactus blockchain!
Please read these guidelines before submitting a pull request or opening an issue.

## Code Guidelines

We strive to maintain clean, readable, and maintainable code in the Pactus blockchain.
Please follow these guidelines when contributing to the project:

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
- Follow the [Go Doc Comments](https://go.dev/doc/comment) guidelines.
- Follow the principles of clean code as outlined in
  Robert C. Martin's "[Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)" book.
- Write tests for new code or changes to existing code, and make sure all tests pass before submitting a pull request.

### Makefile Targets

You can use these commands in the Makefile:

- `make build` compiles the code into executable binaries.
- `make build_gui` compiles the gtk GUI code into executable binary.
- `make devtools` installs required development tools.
- `make fmt` formats the code according to the Go standards.
- `make check` runs checks on the code, including formatting and linting.
- `make test` runs all the tests, including unit tests and system tests.
- `make unit_test` runs only unit tests.
- `make proto` generates [protobuf](https://protobuf.dev/) files.
  Run this target if you have made any changes to the proto buffer files.

### GUI Development

Development of the Pactus Core GUI have some requirements on your machine which you can find a [quick guide about it here](./docs/gtk-gui-development.md).

### Error and Log Messages

Error and log messages should not start with a capital letter (unless it's a proper noun or acronym) and
should not end with punctuation.

All changes on core must contain proper and well-defined unit-tests, also previous tests must be passed as well. 
This codebase used `testify` for unit tests, make sure you follow these guide for tests:

- For panic cases make sure you use `assert.Panics`
- For checking err using `assert.ErrorIs` make sure you pass expected error as second argument.
- For checking equality using `assert.Equal` make sure you pass expected value as the first argument.


> This code guideline must be followed for both contributors and maintainers to review the PRs.

#### Examples

- Correct ✅: "unable to connect to server"
- Incorrect ❌: "Unable to connect to server"
- Incorrect ❌: "unable to connect to server."

### Help Messages

Follow these rules for help messages for CLI commands and flags:

- Help string should not start with a capital letter.
- Help string should not end with punctuation.
- Don't include default value in the help string.
- Include the acceptable range for the flags that accept a range of values.

## Commit Guidelines

Please follow these rules when committing changes to the Pactus blockchain:

- Each commit should represent a single, atomic change to the codebase.
  Avoid making multiple unrelated changes in a single commit.
- Use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for commit messages and
  Pull Request titles.

### Commit type

List of conventional commit [types](https://github.com/commitizen/conventional-commit-types/blob/master/index.json):

| Types    | Description                                                                       |
| -------- | --------------------------------------------------------------------------------- |
| fix      | A big fix                                                                         |
| feat     | A new feature                                                                     |
| docs     | Documentation only changes                                                        |
| test     | Adding missing tests or correcting existing tests                                 |
| build    | Changes that affect the build system or external dependencies                     |
| ci       | Changes to our CI configuration files and scripts                                 |
| perf     | A code change that improves performance                                           |
| refactor | A code change that neither fixes a bug nor adds a feature                         |
| style    | Changes that do not affect the meaning of the code (white-space, formatting, etc) |
| chore    | Other changes that don't modify src or test files                                 |

### Commit Scope

The scope helps specify which part of the code is affected by your commit.
It must be included in the commit message to provide clarity.
Multiple scopes can be used if the changes impact several areas.

Here’s a list of available scopes: [available scopes](./.github/workflows/semantic-pr.yml).

### Commit Description

- Keep the commit message under 50 characters.
- Start the commit message with a lowercase letter and do not end with punctuation.
- Write commit messages in the imperative: "fix bug" not "fixed bug" or "fixes bug".

### Examples

  - Correct ✅: "feat(grpc): sign transaction using wallet client"
  - Correct ✅: "feat(grpc, wallet): sign transaction using wallet client"
  - Incorrect ❌: "feat(gRPC): Sign transaction using wallet client."
  - Incorrect ❌: "feat(grpc): Sign transaction using wallet client."
  - Incorrect ❌: "feat(grpc): signed transaction using wallet client"
  - Incorrect ❌: "sign transaction using wallet client"

## Code of Conduct

This project has adapted the
[Contributor Covenant, version 2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/)
to ensure that our community is welcoming and inclusive for all.
Please read it before contributing to the project.

---

Thank you for your contributions to the Pactus blockchain!
