## Patching Process

To patch a released version, start by creating a dedicated branch for that version, proceed with the required updates, and eventually release the patched branch.

1. **Patch branch**

If this is the first patch for the specific major version, we need to first create a branch for this tag:

```bash
git checkout -b 0.<minor>.x v0.<minor>.0
git log
git push --set-upstream origin 0.<minor>.x
```

and update the patch version inside the [version.go](../version/version.go) file.

Otherwise, switch to the patch branch:

```bash
git checkout origin/0.<minor>.x
```

2. **Updating the branch**

Apply the fixes to the branch. You can use [cherry-pick](https://www.atlassian.com/git/tutorials/cherry-pick) to pick some commits from the main branch and apply them to the patch branch:

```bash
git cherry-pick <commit-id>
git push
```

3. **Creating Environment Variables**

Let's create environment variables for the patch version. For the rest of this document, we will use these environment variables in the commands.

```bash
PRV_VER="0.18.3"
CUR_VER="0.18.4"
NEXT_VER="0.18.5"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
BASE_BRANCH="0.18.x"
```

For the rest of this document, we will use these environment variables in commands.
Keep your terminal open.

4. **Follow the Releasing Document**

Please refer to the [Releasing](./releasing.md) document and follow the steps outlined from Step 4 until the end to complete the patching process.
