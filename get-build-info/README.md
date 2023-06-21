# get-build-info

Determines all metadata for current build and would typically be invoked at beginning of build workflow, within git clone directory and in specific branch being built.

# Inputs

- `work-dir`: Optional working directory within which to execute this action, defaulting to `GITHUB_WORKSPACE` env var. Allows name of project within a monorepo to have its name determined based on `PROJECT` var from its `jen.yaml`. No need to specify this if you specify `project` explicitly.
- `project`: Optional project name override.
- `git-tag-prefix`: Optional prefix for git tags used to version project. In a typical project, that would be simply "v" (the default), eg: "v0.0.1". In a monorepo, that should include the path to sub-project, eg: "path/to/my-project/v".

# Outputs

- `project`: Name of project being built, resolved in this order of precedence:
  - Explicit `project` input passed to this action.
  - `PROJECT` var from `jen.yaml` file in project's root.
  - GitHub repo name derived from `GITHUB_REPOSITORY` env var (eg: `nestoca/hello-world` results in `hello-world`).
  - Current working directory name.
- `version`: Version being built in semver format (eg: "0.0.1" for master or
  "0.1350.0-feat-mo-123-some-feature-6e9580872-2023-05-09_20_17_42" for PRs).
- `releases`: List of releases to promote (defaults to project name), space delimited.
- `git-tag`: Git tag to mark current build in git (eg: "v0.0.1" or
  "path/to/my-project/v0.0.1").
- `docker-tag`: Tag for docker image being built (eg: "0.0.1" for master or
  "0.1350.0-feat-mo-123-some-feature-6e9580872-2023-05-09_20_17_42" for PRs).

# Example

```yaml
jobs:
  job:
    name: Build
    runs-on: ubuntu-latest

    steps:
      - name: Check out
        uses: actions/checkout@v3

      - name: Get build info
        id: info
        uses: nestoca/actions/get-build-info@v1

      - name: Build something with that info
        run: |-
          echo Project: ${{ steps.info.outputs.project }}
          echo Version: ${{ steps.info.outputs.version }}
          echo Releases: ${{ steps.info.outputs.releases }}
          echo Git tag: ${{ steps.info.outputs.git-tag }}
          echo Docker tag: ${{ steps.info.outputs.docker-tag }}
```
