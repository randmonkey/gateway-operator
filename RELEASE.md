# Release Process

To release, [create a new issue](https://github.com/kong/gateway-operator/issues/new/choose) from the "Release" template.

Fill out the issue title and release type, create the issue, and proceed through the release steps, marking them done as you go.

## Release Schedule

During the Technical Preview phase, a technical time-based release of Gateway Operator shall happen monthly on **the last Wednesday of every month**, starting with 0.2.0 on October 26, 2022.

Guideline to pick a version number:
* Look at the set of changes being released. If it has meaningful changes to functionality, API changes, breaking changes, etc. - bump the minor version number (`0.x.y` -> `0.x+1.0`)
* If only bugfixes are being released, bump the patch version number (`0.x.y` -> `0.x.y+1`).
