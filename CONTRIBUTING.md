# Syncing with upstream

This repository is synchronized using [`vendir`](https://github.com/carvel-dev/vendir).

To update the chart to the latest version:

- Execute `APPLICATION=helm/gateway-api make update-chart`
- Update the `appVersion` and `version` fields in `helm/gateway-api/Chart.yaml`
- Add a Changelog entry
- Create a PR to be merged into `main` branch

Once your changes have been tested in CI and are merged into the `main` branch, someone with the correct rights will trigger release creation by pushing a new branch called `main#release#vX.Y.Z`. `X.Y.Z` should match the release of Gateway API we synced to.
