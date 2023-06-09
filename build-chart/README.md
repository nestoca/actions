# build-chart

Builds helm chart and pushes it to nesto's `charts` Google Artifacts Registry repo at:
northamerica-northeast1-docker.pkg.dev/nesto-ci-78a3f2e6

# Inputs

- `name`: Name of chart.
- `version`: Version of chart.
- `key`: Key for service account to use for pushing charts (typically `secrets.GCP_KEY`). Generated by `global` pulumi project, which exposes it as `gitHubActionsSAKeyBase64` output (see: https://github.com/nestoca/infra/blob/47ccfe7e4c34d8fb9cb77a1620a07d8a4acc5099/pulumi/global/components/gcp/projects/cicd.ts#L176).
- `work-dir`: Optional working directory to build chart from (defaults to `.`).
- `registry`: Optional helm chart registry host name (defaults to 'northamerica-northeast1-docker.pkg.dev').

# Outputs

No outputs.

# Example

```yaml
jobs:
  job:
    name: Build image
    runs-on: ubuntu-latest

    steps:
      - name: Check out
        uses: actions/checkout@v3

      - name: Get build info
        id: info
        uses: nestoca/actions/get-build-info@v1

      - name: Build chart
        uses: nestoca/actions/build-chart@v1
        with:
          name: ${{ steps.info.outputs.project }}
          version: ${{ steps.info.outputs.version }}
          key: ${{ secrets.GCP_DOCKER_SA_KEY }}
          work-dir: charts/generic
```
