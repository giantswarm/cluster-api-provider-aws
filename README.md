# Cluster API Provider AWS

This is Giant Swarm's fork. See the upstream [cluster-api-provider-aws README](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/README.md) for official documentation.

## How to work with this repo

We always work on a `release-X.Y` branch, and have no `main` branch.

You create the `release-X.Y` branch like this:

- Create the Git branch from the latest commit in that branch of the upstream repo. That normally means commits from CAPA vX.Y.Z + some backported commits that went into that release branch as well.
- Apply our fork customizations on top. Always compare to our previous `release-X.Y` branch for what was done.
  - `.circleci/config.yml`
  - Change/delete some workflows in `.github/workflows` that our fork shouldn't run, or run differently
  - Delete `.github/PULL_REQUEST_TEMPLATE.md`
  - `CODEOWNERS` file so the right Giant Swarm team is automatically added to PR review requests.
- Test in an MC. We deploy by creating custom Git tags in _this_ repo (see instructions below), and via the chart [cluster-api-provider-aws-app](https://github.com/giantswarm/cluster-api-provider-aws-app/). To provide the YAML manifests, we use GitHub releases, similar to the upstream project. The scripts in `cluster-api-provider-aws-app` convert them into the final manifests.
- Make it the default branch so that PRs in the fork target it

### Test and release

- Do you work and push a Git tag

  ```sh
  # Get latest changes on our release branch
  git checkout release-X.Y
  git pull

  git switch -c my-working-branch

  # Make some changes and commit as usual
  git commit

  # Create and push the Git tag. Fill `BASED_ON_VERSION` with the latest upstream version
  # on which our release branch is based.
  #
  # This triggers the GitHub release action - please continue reading below!
  # The branch needs not be pushed until you want to open a PR. The tag builds
  # the manifest and GitHub release which you need for deploying the app from the
  # `giantswarm/cluster-api-provider-aws-app` repo.
  BASED_ON_VERSION=vX.Y.Z; [[ $BASED_ON_VERSION =~ ^v[0-9] ]] && git tag "${BASED_ON_VERSION}-gs-$(git rev-parse --short HEAD)" && git push origin "${BASED_ON_VERSION}-gs-$(git rev-parse --short HEAD)"
  ```

- Check that the [CircleCI pipeline](https://app.circleci.com/pipelines/github/giantswarm/cluster-api-provider-aws) succeeds for the desired Git tag in order to produce images. If the tag build fails, fix it.
- Edit [that draft GitHub release](https://github.com/giantswarm/cluster-api-provider-aws/releases) and turn it from draft to released. This makes the release's manifest files available on the internet, as used in [cluster-api-provider-aws-app](https://github.com/giantswarm/cluster-api-provider-aws-app). If you see no GitHub release after a few minutes, check that the [GitHub release action](https://github.com/giantswarm/cluster-api-provider-aws/actions) succeeded for your tag.
- Test the changes in the app
  - Replace `.tag` in [cluster-api-provider-aws-app's `values.yaml`](https://github.com/giantswarm/cluster-api-provider-aws-app/blob/master/helm/cluster-api-provider-aws/values.yaml) with the new tag `vX.Y.Z-gs-...`.
  - Run `cd cluster-api-provider-aws-app && make generate` to update manifests
  - Commit and push your working branch for `cluster-api-provider-aws-app` to trigger CircleCI pipeline
  - Install and test the app thoroughly on a management cluster. Continue with the next step only once you're confident.
- Open PR for `cluster-api-provider-aws` fork (your working branch)
  - If you merged an upstream release tag, we should target our `release-X.Y` branch with the PR.
  - On the other hand, if you implemented something else which is not in upstream yet, we should target `upstream/main` so that it first lands in the upstream project, officially approved, tested and released. Afterwards, you would repeat this whole procedure and merge the release that includes your fix. For a quick in-house hotfix, you can alternatively do a quicker PR targeted against our `release-X.Y` branch.
- Also open PR for `cluster-api-provider-aws-app` change
- Once merged, it should automatically be [deployed via management-cluster-bases Flux resources](https://github.com/giantswarm/management-cluster-bases/blob/main/bases/collections/capa/base/cluster-api-provider-aws.yaml)
