---
name: "🚀 Release"
about: 'Tracking a new release of the Kong Kubernetes Gateway Operator'
title: ''
labels: ''
assignees: ''

---

## Steps

This release pipeline is under continous improvement. If you encounter any problem, please refer to the [troubleshooting](#troubleshooting) section of this document. If the troubleshooting section does not contain the answer to the problem you encountered, please create an issue to improve either the pipeline (if the problem is a bug), or this document (if the problem is caused by human error).

- [ ] Check [default versions](#verify-default-hardcoded-versions) of images (see below).
- [ ] Check the `CHANGELOG.md` and update it with the new version number. Make sure the log is up to date.
- [ ] Check the Kong incubator [website][kongincubator] ([KGO page source][kongincubator-kgo-project]) and update it with the released documentation.
- [ ] Check the existing [releases][releases] and determine the next version number.
- [ ] From [GitHub release action][release-action], start a new workflow run with the `release` input set to the release tag (e.g. `v0.1.0`).
- [ ] Wait for the workflow to complete.
- [ ] Ensure CI created a [KGO docs][kgo-docs-prs] PR. Review and merge it. 
- [ ] The CI should create a PR in the [Gateway Operator][kgo-prs] repo that syncs the release branch to the `main` branch. Merge it.
- [ ] After the PR is merged, a new release should be created automatically. Check the [releases][releases] page.
- [ ] Submit the operator to [external hubs](#submit-to-external-hubs) (see below).

## Verify default hardcoded versions

The package [internal/consts][consts-pkg] contains a list of default versions for the operator. These versions should be updated to match the new release. The example consts to look for:

- `DefaultDataPlaneTag`
- `DefaultControlPlaneTag`
- `WebhookCertificateConfigBaseImage`

Also, the Makefile contains hardcoded information that needs to be updated:

- `CHANNELS` - the channels the operator is available on the OpenShift Operator Hub. For the technical preview we're using `alpha` channels. Please refer to the [OLM docs][olm-channels] for more information.
- `OPENSHIFT_SUPPORTED_VERSIONS` - the supported versions of OpenShift.

When the changes of the above versions are ready, make sure you run `make generate manifests bundle.regular`.

## Submit to external hubs

- [ ] Submit the operator to the [OperatorHub](#operatorhub-community-operators-steps) and wait for it to be published (follow instruction below).

### OperatorHub Community Operators steps

[Operator Hub Community Operators][operator-hub-community]

- [ ] The PR to the [community operators][operator-hub-community] repository should be created by the CI. Please check that the PR exists in the [community operators][operator-hub-community] repository.

## Troubleshooting

### The release needs to be started again with the same tag

If the release workflow needs to be started again with the same input version, the release branch needs to be deleted. The release branch is created by the CI and it's named `release/v<version>`. For example, if the release version is `v0.1.0`, the release branch will be `release/v0.1.0`.

It's only safe to start a release workflow with the version that was previously used if:

- The release PR to the gateway-operator repo is not merged yet.
- The external hub PRs are not merged yet.
- The tag that matches the release version does not exist.

Otherwise, if the above conditions are not meet, the release cannot be restarted. A new release needs to be started with an input version that would be next in semantic versioning.

Steps:

1. Delete the `release/v<version>` branch.
2. Delete the PR created by a release workflow.
3. Update the repository with the correct changes.
4. Start a new release workflow run.

### OperatorHub Community Operators PR failed

When the release workflow is restarted with the same input version, the OperatorHub Community Operators PR might fail. This is because the PR already exists in the [community operators][operator-hub-community] repository. The PR needs to be closed manually.

#### Option 1: Manually fix the the PR

Steps:

1. Checkout the [community operators fork][community-operators-fork] repository.
2. Checkout the branch named `kong-gateway-operator-<version>`.
3. Fix the PR.
4. Commit the changes. The commits need to have a `signed-off-by` clause. (`git commit -s`)
5. Push the changes to the [community operators fork][community-operators-fork].

#### Option 2: Re-run the workflow

It's only safe to do so if:

- The release PR to the gateway-operator repo is not merged yet.
- The external hub PRs are not merged yet.
- The tag that matches the release version does not exist.

Otherwise, if the above condition are not meet, the release cannot be restarted. Instead, use the [Option 1](#option-1-manually-fix-the-the-pr).

Steps:

1. Login to GitHub as `team-k8s-bot`
2. Delete the PR together with the branch from the [community operators fork][community-operators-fork] repository.
3. Delete the release from the [releases][releases] page.
4. Start a new release workflow run.

[releases]: https://github.com/Kong/gateway-operator/releases
[release-action]: https://github.com/Kong/gateway-operator/actions/workflows/release.yaml
[community-operators-fork]: https://github.com/kong/k8s-operatorhub-community-operators
[consts-pkg]: https://github.com/Kong/gateway-operator/blob/main/internal/consts/consts.go
[olm-channels]: https://olm.operatorframework.io/docs/best-practices/channel-naming/
[operator-hub-community]: https://github.com/k8s-operatorhub/community-operators
[kongincubator]: https://incubator.konghq.com/p/gateway-operator
[kongincubator-kgo-project]: https://github.com/Kong/kong-incubator/blob/main/src/_projects/gateway-operator.md
[kgo-docs-prs]: https://github.com/Kong/gateway-operator-docs/pulls
[kgo-prs]: https://github.com/Kong/gateway-operator/pulls
