# Release Candidate

Release candidates are used to test the latest version before a final release.
They are **not intended to be stored permanently** and will be removed from GitHub Releases after testing.
To create a release candidate, follow the structure below:

## 1. Set Environment Variables

Create environment variables for the release version, which will be used in subsequent commands throughout this document.
Keep your terminal open for further steps.
Update `X`, `Y`, `Z`, and `N` to reflect the the release candidate (e.g., `1.7.0-rc1`, `1.7.0-rc2`, etc.).

```bash
CUR_VER="X.Y.Z-rcN"
TAG_NAME="v${CUR_VER}"
TAG_MSG="Version ${CUR_VER}"
```

## 2. Update the Version

Set `rcX` in the `Meta` field and update the `Alias` in [version.go](../version/version.go).

## 3. Make a Release Candidate Commit

Create a new commit for the release candidate:

```bash
git commit -a -m "chore(release): releasing version ${CUR_VER}"
```

## 4. Tag the Release

Create a signed Git tag:

```bash
git tag -s -a ${TAG_NAME} -m "${TAG_MSG}"
```

Verify the tag:

```bash
git show ${TAG_NAME}
```

## 5. Push the Tag

Push the tag to the remote repository:

```bash
git push origin ${TAG_NAME}
```

Pushing the tag will automatically trigger the creation of a release candidate and build the binaries.
