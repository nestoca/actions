# build-image

Builds docker image and pushes it to one of nesto's GCP registries:

- Google Container Registry: gcr.io/nesto-ci-78a3f2e6
- Google Artifacts Registry: northamerica-northeast1-docker.pkg.dev/nesto-ci-78a3f2e6

# Inputs

- `name`: Name of image to build and push.
- `tags`: Tags to put on image.
- `key`: Key for service account to use for pushing images (typically `secrets.GCP_KEY`). Generated by `global` pulumi project, which exposes it as `gitHubActionsSAKeyBase64` output (see: https://github.com/nestoca/infra/blob/47ccfe7e4c34d8fb9cb77a1620a07d8a4acc5099/pulumi/global/components/gcp/projects/cicd.ts#L176).
- `context`: Optional context directory to build image from (defaults to `.`).
- `target`: Optional registry/repo to push image to, either `default` (the default) or `actions`.

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

      - name: Build image
        uses: nestoca/actions/build-image@v1
        with:
          name: ${{ steps.info.outputs.project }}
          tags: ${{ steps.info.outputs.docker-tag }},latest
          key: ${{ secrets.GCP_KEY }}
```