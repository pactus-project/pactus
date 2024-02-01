# Patching Process

To patch a released version, start by creating a dedicated branch for the patch version,
proceed with the required updates, and eventually release the patched branch.

## 1. Preparing Your Environment

Before proceeding with the patching process,
ensure that your `origin` remote is set to `git@github.com:pactus-project/pactus.git` and not your local fork.

## 2. Create a Patch Branch

If this is the first patch for a specific major version, you'll need to create a branch for this tag.
Replace `<minor>` with the appropriate minor version number:

```bash
git checkout -b 0.<minor>.x v0.<minor>.0
git log
git push --set-upstream origin 0.<minor>.x
```

Don't forget to update the patch version inside the [version.go](../version/version.go) file.

If you're not creating a new patch branch, switch to the existing patch branch:

```bash
git checkout 0.<minor>.x
git pull
```

## 3. Apply Fixes

Now, apply the necessary fixes to the patch branch.
You can use [cherry-pick](https://www.atlassian.com/git/tutorials/cherry-pick) to
select specific commits from the main branch and apply them to the patch branch:

```bash
git cherry-pick <commit-id>
git push
```

## 4. Set Environment Variables

Reopen this document within the branch version and
create environment variables for the release version, which will be used in subsequent commands throughout this document.
Keep your terminal open for further steps.

```bash
PRV_VER="1.0.0"
CUR_VER="1.0.1"
NEXT_VER="1.0.2"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
BASE_BRANCH="0.21.x"
```

## 5. Follow the Releasing Document

Refer to the [Releasing](./releasing.md) document and follow the steps outlined from Step 5 until the end.
This document will provide you with the necessary guidance to successfully release the patched branch.

Ensure that your terminal remains open throughout the process for seamless execution of the required commands.
