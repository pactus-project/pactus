## Patching Process

To patch a released version, start by creating a dedicated branch for that version, proceed with the required updates, and eventually release the patched branch.

1. **Creating Environment Variables**

Let's create environment variables for the patch version. For the rest of this document, we will use these environment variables in the commands.

```bash
PRV_VER="0.18.0"
CUR_VER="0.18.0"
NEXT_VER="0.18.1"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
BASE_BRANCH="0.18.x"
BASE_VER="0.18.0"
```

For the rest of this document, we will use these environment variables in commands.
Keep your terminal open.

2. **Creating a new branch**

If this is the first patch for the specific major version, we need to first create a branch for this tag:

```bash
git checkout -b ${BASE_BRANCH} v${BASE_VER}
git branch --set-upstream-to=origin/${BASE_BRANCH}
```

and update the patch version inside the [version.go](../version/version.go) file.

Otherwise, switch to the patch branch:

```bash
git checkout ${BASE_BRANCH}
```

3. **Updating the branch**

Apply the fixes to the branch. You can use [cherry-pick](https://www.atlassian.com/git/tutorials/cherry-pick) to pick some commits from the main branch and apply them to the patch branch:

```bash
git cherry-pick <commit-id>
```

4. **Follow the [Releasing](./releasing.md) Document**

Please refer to the [Releasing](./releasing.md) document and follow the steps outlined from Step 4 until the end to complete the patching process.
