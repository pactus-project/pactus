# Contributing

Thank you for considering contributing to Pactus blockchain!
Please read these guidelines before submitting a pull request or opening an issue.

## Code guidelines

We strive to maintain clean, readable, and maintainable code in Pactus blockchain.
Please follow these guidelines when contributing code to the project:

- Code should follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
- Documentation should follow the [Go Doc Comments](https://go.dev/doc/comment) format.
- Use clear and descriptive variable, function, and method names. Avoid using abbreviations or acronyms unless they are well-known and widely understood.
- Write code that is modular and reusable, and follows the [DRY (Don't Repeat Yourself)](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) principle.
- Write tests for new code or changes to existing code, and make sure all tests pass before submitting a pull request.
- Error strings and log messages SHOULD NOT be capitalized (unless beginning with proper nouns or acronyms) and
 SHOULD NOT end with punctuation. Examples:
  * Correct: "unable to connect to server"
  * Incorrect: "Unable to connect to server"
  * Incorrect: "unable to connect to server."

Before submitting a pull request, please run the following commands to ensure there are no issues:

- `make fmt` will format the code according to the Go standards.
- `make check` will run various checks on the code, including formatting and linting.
- `make test`  will run the tests to ensure that all functionality is working as intended.

## Commit guidelines

Please follow these guidelines when committing changes to Pactus blockchain:

- Each commit should represent a single, atomic change to the codebase.
  Avoid making multiple unrelated changes in a single commit.
- Use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for commit messages and Pull Request titles.

Common prefixes:

| Types    | Description                                                   |
| -------- | ------------------------------------------------------------- |
| fix      | A big fix                                                     |
| feat     | A new feature                                                 |
| docs     | Documentation changes                                         |
| test     | Adding missing tests or correcting existing tests             |
| build    | Changes that affect the build system or external dependencies |
| ci       | Changes to our CI configuration files and scripts             |
| perf     | A code change that improves performance                       |
| refactor | A code change that neither fixes a bug nor adds a feature     |

## Code of Conduct

This project has adapted the [Contributor Covenant, version 2.1](https://www.contributor-covenant.org/version/2/1/code_of_conduct/)
to ensure that our community is welcoming and inclusive for all.
Please read it before contributing to the project.

---

Thank you for your contributions to Pactus blockchain!
