# Releasing

Before releasing and publishing the Pactus software, there are a few important steps that need to be followed.
Please follow the instructions below:

1. Get the latest code

```bash
git checkout main
git pull
```

2. Updating Windows DLLS

To ensure that the GUI can find the required dependency DLLs in Windows, we may need to update them.
Follow these commands in the project's root directory, using[MSYS2](https://www.msys2.org/):

```bash
git pull
pacman -Suyyy
.github/releasers/releaser_gui_windows.sh
```

Wait for the build to finish. If everything is okay, proceed to the next step.
Otherwise, update the dependency DLLs inside `.github/releasers/releaser_gui_windows.sh` and
run the above command again.

3. Check the config

Double-check the config files to ensure they are up to date.

4. Creating Environment Variables

Let's create environment variables for the release version.
For the rest of this document, we will use these environment variables in the commands.

```bash
VERSION="0.11.0"
TAG_NAME="v${VERSION}"
TAG_MESSAGE="Version ${VERSION}"
```

For the rest of this document, we will use these environment variables in commands.

5. Update Changelog

Use [Commitizen](https://github.com/commitizen-tools/commitizen) to update the CHANGELOG.
Run the following command:

```bash
cz changelog --incremental --unreleased-version $VERSION
```

Sometimes you may need to amend the changelog manually.
Create a comparison link for the changelog header, like:

```text
## [0.11.0](https://github.com/pactus-project/pactus/compare/v0.10.0...v0.11.0)
```

6. Create release PR

Create a new PR against the `main` branch:

```bash
git checkout -b releasing_$VERSION
git commit -a -m "chore: Releasing version $VERSION"
git push origin HEAD
```

Wait for the PR to be approved and merged into the `main` branch.

7. Tagging

Create a git tag:

```bash
git checkout main
git pull
git tag -s -a $TAG_NAME -m $TAG_MESSAGE
```

check the tag info:

```bash
git show $TAG_NAME
```

8. Push the tag

Now you can push the tag to the repository:

```bash
git push origin $TAG_NAME
```

Pushing the tag will automatically create a release tag and build the binaries.

9. Bumping version

Update the version inside the `version/version.go` to `0.12.0`
Also update the version inside this document in step 3 and 4 to `0.12.0`

Create a new PR against `main` branch:

```bash
git checkout -b bumping_0.12.0
git commit -a -m "chore: bumping version to 0.12.0"
git push origin HEAD
```

Wait for the PR to be approved and merged into the `main` branch.

10. Update  the website

Create a new announcement post in the [blog](https://pactus.org/blog/) and
update the [Road Map](https://pactus.org/about/roadmap/).
Additionally, draft a new release on the
[GitHub Releases](https://github.com/pactus-project/pactus/releases) page.

11. Celebrate 🎉
