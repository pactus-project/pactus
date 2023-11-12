# Releasing

Before releasing and publishing the Pactus software, there are a few important steps that need to be followed.
Please follow the instructions below:

1. **Get the latest code**

```bash
git checkout main
git pull
```

2. **Updating Windows DLLS**

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

3. **Creating Environment Variables**

Let's create environment variables for the release version.
For the rest of this document, we will use these environment variables in the commands.

```bash
PRV_VER="0.16.0"
CUR_VER="0.17.0"
NEXT_VER="0.18.0"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
BASE_BRANCH="main"
```

For the rest of this document, we will use these environment variables in commands.
Keep your terminal open.

4. **Update the version**

Remove `beta` from the `meta` field in [version.go](../version/version.go).
Also, double-check the config files to ensure they are up to date.

5. **Update Changelog**

Use [Commitizen](https://github.com/commitizen-tools/commitizen) to update the CHANGELOG.
Run the following command:

```bash
cz changelog --incremental --unreleased-version ${TAG_NAME}
```

Sometimes you may need to amend the changelog manually.
Then, add links to the CHANGELOG:

```bash
sed -E -i "s/## v${CUR_VER} /## [${CUR_VER}](https:\/\/github.com\/pactus-project\/pactus\/compare\/v${PRV_VER}...v${CUR_VER}) /g" CHANGELOG.md
sed -E -i 's/\(#([0-9]+)\)/([#\1](https:\/\/github.com\/pactus-project\/pactus\/pull\/\1))/g' CHANGELOG.md
```

6. **Create release PR**

Create a new PR against the base branch.
We use [GiyhUb CLI](https://github.com/cli/cli/) to create the PR, but you can create it manually.

```bash
git checkout -b releasing_${CUR_VER}
git commit -a -m "chore: releasing version ${CUR_VER}"
git push origin HEAD
gh pr create --title "chore: releasing version ${CUR_VER}" --body "Releasing version ${CUR_VER}" --base ${BASE_BRANCH}
```

Wait for the PR to be approved and merged into the base branch.

7. **Tagging**

Create a git tag:

```bash
git checkout ${BASE_BRANCH}
git pull
git tag -s -a $TAG_NAME -m $TAG_MSG
```

check the tag info:

```bash
git show $TAG_NAME
```

8. **Push the tag**

Now you can push the tag to the repository:

```bash
git push origin $TAG_NAME
```

Pushing the tag will automatically create a release tag and build the binaries.

9. **Bumping version**

Update the version inside the [version.go](../version/version.go) and add `beta` to the `meta` field.
Update [patching](./patching.md) docuemnt.
If this is a majore release, update the version inside this document in step 3.

Create a new PR against base branch:

```bash
git checkout -b bumping_${NEXT_VER}
git commit -a -m "chore: bumping version to ${NEXT_VER}"
git push origin HEAD
gh pr create --title "chore: bumping version to ${NEXT_VER}" --body "Bumping version to ${NEXT_VER}" --base ${BASE_BRANCH}
```

Wait for the PR to be approved and merged into the base branch.

10. **Update  the website**

Create a new announcement post in the [blog](https://pactus.org/blog/) and
update the [Road Map](https://pactus.org/about/roadmap/) and
[Download](https://pactus.org/download/) pages.
Additionally, draft a new release on the
[GitHub Releases](https://github.com/pactus-project/pactus/releases) page.

11. **Celebrate 🎉**

Before celebrating, ensure that the release has been tested and that all documentation is up to date
