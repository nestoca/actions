# get-build-info

That action allows to determine all the metadata for current build and would typically be invoked at beginning of build workflow, within git clone directory and in specific branch being built.

# Inputs

- `project`: Optional project name override, defaulting to PROJECT var in jen.yaml if present, otherwise uses current working directory name.
- `git-tag-prefix`: Optional prefix for git tags used to version project. In a typical project, that would be simply "v" (the default), eg: "v0.0.1". In a monorepo, that should include the path to sub-project, eg: "path/to/my-project/v".

# Outputs

- `project`: Name of project being built.
- `version`: Version being built in semver format (eg: "0.0.1" for master or
  "0.1350.0-feat-mo-123-some-feature-6e9580872-2023-05-09_20_17_42" for PRs).
- `git-tag`: Git tag to mark current build in git (eg: "v0.0.1" or
  "path/to/my-project/v0.0.1").
- `docker-tag`: Tag for docker image being built (eg: "0.0.1" for master or
  "0.1350.0-feat-mo-123-some-feature-6e9580872-2023-05-09_20_17_42" for PRs).

# Example

```yaml
jobs:
  job:
    name: Build my-project
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3
      - id: info
        uses: nestoca/actions/get-build-info@master
        with:
          project: my-project
          git-tag-prefix: v
      - run: |
          echo Project: ${{ steps.info.outputs.project }}
          echo Version: ${{ steps.info.outputs.version }}
          echo Git-Tag: ${{ steps.info.outputs.git-tag }}
          echo Docker-Tag: ${{ steps.info.outputs.docker-tag }}
```
