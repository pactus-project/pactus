# Patching Process

To patch a released version, start by creating a dedicated branch for the patch version,
proceed with the required updates, and eventually release the patched branch.

## 1. Preparing Your Environment

Before proceeding with the patching process,
ensure that your `origin` remote is set to `git@github.com:pactus-project/pactus.git` and not your local fork.

## 2. Set Environment Variables

Create environment variables for the patch version, which will be used in subsequent commands throughout this document.
Keep your terminal open for further steps.

```bash
PRV_VER="1.7.0"
CUR_VER="1.7.1"
NEXT_VER="1.7.2"
BASE_BRANCH="1.7.x"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
```

## 3. Patch Branch

### Create the Patch branch

If this is the first patch for a specific major version, you'll need to create a branch for this tag.

```bash
git checkout -b ${BASE_BRANCH} v${PRV_VER}
git log
git push --set-upstream origin ${BASE_BRANCH}
```

Update the patch version inside the [version.go](../version/version.go),
clear Alias and set "beta" to the Meta.

Create a new PR against the patch branch:

```bash
git checkout -b bumping_${CUR_VER}
git commit -a -m "chore: bumping version to ${CUR_VER}"
git push origin HEAD
gh pr create --title "chore: bumping version to ${CUR_VER}" --body "Bumping version to ${CUR_VER}" --base ${BASE_BRANCH}
```

As an example check this [Pull Request](https://github.com/pactus-project/pactus/pull/1367).
Wait for the PR to be approved and merged into the patch branch.

### Switch the Patch branch

If you're not creating a new patch branch, switch to the existing patch branch:

```bash
git checkout ${BASE_BRANCH}
git pull
```

## 4. Apply Fixes

Now, apply the necessary fixes to the patch branch.
You can use [cherry-pick](https://www.atlassian.com/git/tutorials/cherry-pick) to
select specific commits from the main branch and apply them to the patch branch:

```bash
git cherry-pick <commit-id>
git push
```

## 5. Follow the Releasing Document

Refer to the [Releasing](./releasing.md) document and follow the steps outlined from Step 5 until the end.
This document will provide you with the necessary guidance to successfully release the patched branch.

Ensure that your terminal remains open throughout the process for seamless execution of the required commands.
