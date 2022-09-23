---
name: "🚀 Release"
about: 'Tracking a new release of the Kong Kubernetes Gateway Operator'
title: ''
labels: ''
assignees: ''

---

## Steps

- [ ] Check [default versions](#verify-default-hardcoded-versions) of images.
- [ ] Check the `CHANGELOG.md` and update it with the new version number. Make sure the log is up to date.
- [ ] Check the existing [releases][releases] and determine the next version number.
- [ ] From [GitHub release action][release-action], start a new workflow run with the `release` input set to the release tag (e.g. `v0.1.0`).
- [ ] Wait for the workflow to complete.
- [ ] Check the release notes from [releases][releases] and make sure everything is OK.
- [ ] Submit the operator to [external hubs](#submit-to-external-hubs).

## Verify default hardcoded versions

The package [internal/consts][consts-pkg] contains a list of default versions for the operator. These versions should be updated to match the new release. The example consts to look for:

- `DefaultDataPlaneTag`
- `DefaultControlPlaneTag`
- `WebhookCertificateConfigBaseImage`

Also, the Makefile contains hardcoded information that needs to be updated:

- `CHANNELS` - the channels the operator is available on the OpenShift Operator Hub. For the technical preview we're using `alpha` channels. Please refer to the [OLM docs][olm-channels] for more information.
- `OPENSHIFT_SUPPORTED_VERSIONS` - the supported versions of OpenShift.

The Redhat Certified Operators bundles has a env vars file that needs to be updated `config/redhat-certified/manager_config_patch.yaml`:
- `RELATED_IMAGE_KONG` - the Kong data plane image.
- `RELATED_IMAGE_KONG_CONTROLLER` - the Kong control plane image.
- `RELATED_IMAGE_CERTIFICATE_CONFIG` - the image used to generate the webhook certificates.


## Submit to external hubs

- [ ] Submit the operator to the [OperatorHub](#operatorhub-community-operators-steps) and wait for it to be published.
- [ ] Submit the operator to the Red Hat [Certified Operators](#red-hat-certified-operators-steps) and wait for it to be published.

### OperatorHub Community Operators steps

<!-- TODO: This section needs more to be considered done -->
[Operator Hub Community Operators][operator-hub-community]


### Red Hat Certified Operators Steps

Please refer to the [Operator Certification Guide][operator-certification-pipeline] for more information .

- [ ] The release should created a new branch in [the redhat certified operators repository fork][certified-operators-fork] named `kong-gateway-operator-<version>`.
- [ ] Verify that the [the redhat certified operators repository fork][certified-operators-fork] has the new branch and generated bundle is present in subdirectory `operators/kong-gateway-operator/<version>`.
- [ ] Start a tekton pipeline on the OpenShift cluster to test the operator and issue a pull request to the upstream repo

```console
VERSION=<version> # e.g. v0.1.0
BRANCH=kong-gateway-operator-$VERSION

tkn pipeline start operator-ci-pipeline \
  --use-param-defaults \
  --param git_repo_url=git@github.com:Kong/redhat-certified-operators.git \
  --param git_branch=$BRANCH \
  --param bundle_path=operators/kong-gateway-operator/$VERSION \
  --param env=prod \
  --param upstream_repo_name=redhat-openshift-ecosystem/certified-operators \
  --param git_username=team-k8s-bot \
  --param git_email="team-k8s+github-bot@konghq.com" \
  --param pin_digests=true \
  --workspace name=ssh-dir,secret=github-ssh-credentials \
  --workspace name=pipeline,volumeClaimTemplateFile=templates/workspace-template.yml \
  --workspace name=registry-credentials,secret=registry-redhat-dockerconfig-secret \
  --param submit=true
```

The `templates/workspace-template.yml` comes from the [Red Hat ISV Operator Certification Pipelines][operator-pipelines repository]. More information about the pipeline can be found in the [Operator Certification Guide][operator-certification-pipeline]. The content of the file is:

```console
cat templates/workspace-template.yml
---
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```


[releases]: https://github.com/Kong/gateway-operator/releases
[release-action]: https://github.com/Kong/gateway-operator/actions/workflows/release.yaml
[certified-operators-fork]: https://github.com/Kong/redhat-certified-operators/
[certified-operators]: https://github.com/redhat-openshift-ecosystem/certified-operators
[operator-pipelines]: https://github.com/redhat-openshift-ecosystem/operator-pipelines
[operator-certification-pipeline]: https://github.com/Kong/team-k8s/blob/main/docs/operator_certification_pipeline.md
[consts-pkg]: https://github.com/Kong/gateway-operator/blob/main/internal/consts/consts.go
[olm-channels]: https://olm.operatorframework.io/docs/best-practices/channel-naming/
[operator-hub-community]: https://github.com/k8s-operatorhub/community-operators
