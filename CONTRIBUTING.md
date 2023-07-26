# Contributing

Thank you for considering contributing to Pactus blockchain!
Please read these guidelines before submitting a pull request or opening an issue.

## Code guidelines

We strive to maintain clean, readable, and maintainable code in Pactus blockchain.
Please follow these guidelines when contributing code to the project:

- Code should follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
- Documentation should follow the [Go Doc Comments](https://go.dev/doc/comment) format.
- Follow the principles of clean code as outlined in Robert C. Martin's "[Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)" book. Here you can find a [summary](https://gist.github.com/wojteklu/73c6914cc446146b8b533c0988cf8d29) of the book.
- Write tests for new code or changes to existing code, and make sure all tests pass before submitting a pull request.
- Error strings and log messages SHOULD NOT be capitalized (unless beginning with proper nouns or acronyms) and
 SHOULD NOT end with punctuation. Examples:
  * Correct: "unable to connect to server"
  * Incorrect: "Unable to connect to server"
  * Incorrect: "unable to connect to server."

The following commands are available in the Makefile:

- `make devtools` installs required development tools.
- `make fmt` formats the code according to the Go standards.
- `make check` runs various checks on the code, including formatting and linting.
- `make test` runs the tests to ensure that all functionality is working as intended.
- `make proto` regenerates the corresponding code if you have made any changes to the proto buffer files.

## Commit guidelines

Please follow these guidelines when committing changes to Pactus blockchain:

- Each commit should represent a single, atomic change to the codebase.
  Avoid making multiple unrelated changes in a single commit.
- Use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for commit messages and Pull Request titles.

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

## Code of Conduct

This project has adapted the [Contributor Covenant, version 2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/)
to ensure that our community is welcoming and inclusive for all.
Please read it before contributing to the project.

---

Thank you for your contributions to Pactus blockchain!
