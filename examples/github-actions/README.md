# aactl as builder in GitHub Actions (GHA)

In addition to being used as a CLI, `aactl` can also be used as a github action.

## inputs

* `project` - (required) GCP Project ID
* `source` - (required) Full image path with tag or digest
* `file` - (required) Path to the vulnerability file

## usage

Below example, shows how to import vulnerabilities from previously generated report.


```yaml
      - name: 'Run aactl'
        uses: docker://gcr.io/cloud-builders/aactl:latest
        with:
          args: vuln --project ${{ env.PROJECT_ID }} --source ${{ env.IMAGE_ID }} --file ${{ steps.scan.outputs.output }}
```

> Fully working example can be found in [on-push.yaml](on-push.yaml).
