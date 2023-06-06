# promote-build

Promotes build version in given releases of target environment using a certain CD mode.
Typically called at end of build workflow to auto-promote new build to some environment.

# Inputs

- `version`: Version of build to promote.
- `releases`: Space-delimited list of releases to promote (in most cases, will be just the project name).
- `env`: Target environment to promote to (defaults to `qa`).
- `via`: Optional promotion path to use, either "legacy" (using `*-infrastucture` repos and Cloud Build), "codefresh" (using `infra` monorepo, Codefresh and helmfile, the default), or "joy" (using `catalog` repo, Joy and ArgoCD). Only "codefresh" is supported for now.
- `token`: GitHub token to use to check out infra repo.

# Outputs

No outputs.

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

      # Build something here...

      - name: Promote build to qa environment via Codefresh CD
        with:
          version: ${{ steps.info.outputs.version }}
          releases: ${{ steps.info.outputs.releases }}
          env: qa
          via: codefresh
          token: $${{ secrets.GH_TOKEN }}
```
