# Release Process

To ensure a successful release and publication of the Pactus software, it is essential to follow these key steps.
Please carefully follow the instructions provided below:

## 1. Preparing Your Environment

Before proceeding with the releasing process,
ensure that your `origin` remote is set to `git@github.com:pactus-project/pactus.git` and not your local fork.

## 2. Fetch the Latest Code

Ensure that your local repository is up-to-date with the Pactus main repository:

```bash
git checkout main
git pull
```

## 3. Update Windows DLLs

To ensure that the GUI can locate the required dependency DLLs on Windows,
you may need to update them for the [Windows installer](../.github/releasers/releaser_gui_windows.sh).
Make sure you have access to a Windows OS and follow these steps in the project's root directory using [MSYS2](https://www.msys2.org/):

```bash
git pull
pacman -Suyyy
.github/releasers/releaser_gui_windows.sh
```

Wait for the build to finish. If everything is successful, proceed to the next step.
If not, update the dependency DLLs inside `.github/releasers/releaser_gui_windows.sh` and rerun the command.

## 4. Set Environment Variables

Create environment variables for the release version, which will be used in subsequent commands throughout this document.
Keep your terminal open for further steps.

```bash
PRV_VER="1.0.0"
CUR_VER="1.1.0"
NEXT_VER="1.2.0"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
BASE_BRANCH="main"
```

## 5. Update the Version

Remove the `beta` tag from the `meta` field in [version.go](../version/version.go).
Also, double-check the [config.go](../config/config.go) files to ensure they are up-to-date.

## 6. Update Changelog

Use [Commitizen](https://github.com/commitizen-tools/commitizen) to update the CHANGELOG. Execute the following command:

```bash
cz changelog --incremental --unreleased-version ${TAG_NAME}
perl -i -pe "s/## v${CUR_VER} /## [${CUR_VER}](https:\/\/github.com\/pactus-project\/pactus\/compare\/v${PRV_VER}...v${CUR_VER}) /g" CHANGELOG.md
perl -i -pe "s/\(#([0-9]+)\)/([#\1](https:\/\/github.com\/pactus-project\/pactus\/pull\/\1))/g" CHANGELOG.md
```

Occasionally, you may need to make manual updates to the [CHANGELOG](../CHANGELOG.md).

## 7. Create a Release PR

Generate a new PR against the base branch.
It's better to use [GitHub CLI](https://github.com/cli/cli/) to create the PR, but manual creation is also an option.

```bash
git checkout -b releasing_${CUR_VER}
git commit -a -m "chore: releasing version ${CUR_VER}"
git push origin HEAD
gh pr create --title "chore: releasing version ${CUR_VER}" --body "Releasing version ${CUR_VER}" --base ${BASE_BRANCH}
```

Wait for the PR to be approved and merged into the main branch.

## 8. Tagging the Release

Create a Git tag and sign it using your [GPG key](https://docs.github.com/en/authentication/managing-commit-signature-verification/about-commit-signature-verification) with the following commands:

```bash
git checkout ${BASE_BRANCH}
git pull
git tag -s -a $TAG_NAME -m $TAG_MSG
```

Inspect the tag information:

```bash
git show ${TAG_NAME}
```

## 9. Push the Tag

Now, push the tag to the repository:

```bash
git push origin ${TAG_NAME}
```

Pushing the tag will automatically create a release tag and build the binaries.

## 10. Bump the Version

Update the version inside [version.go](../version/version.go) and add `beta` to the `meta` field.
Additionally, update version in the [patching](./patching.md) document.
If this is a major release, update the version inside this document in step 3.

Create a new PR against the base branch:

```bash
git checkout -b bumping_${NEXT_VER}
git commit -a -m "chore: bumping version to ${NEXT_VER}"
git push origin HEAD
gh pr create --title "chore: bumping version to ${NEXT_VER}" --body "Bumping version to ${NEXT_VER}" --base ${BASE_BRANCH}
```

Wait for the PR to be approved and merged into the main branch.

## 11. Update the Website

Create a new announcement post on the
[blog](https://pactus.org/blog/) and update the
[Road Map](https://pactus.org/about/roadmap/) and
[Download](https://pactus.org/download/) pages.
Additionally, draft a new release on the
[GitHub Releases](https://github.com/pactus-project/pactus/releases) page.

## 12. Celebrate 🎉

Before celebrating, ensure that the release has been tested and that all documentation is up to date
