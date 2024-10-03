# Changelog

## [2.7.1](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.7.0...v2.7.1) (2024-10-03)


### Bug Fixes

* fix the values in the feed manifest. ([be509e7](https://github.com/ivan-poltavskiy/news-aggregator/commit/be509e7b19833af11356f024a3fecc3a27d28f93))
* fix the values.yaml and add step for installing certmanager to the install-chart task. ([04e83a5](https://github.com/ivan-poltavskiy/news-aggregator/commit/04e83a5ad75123f7111f27b10a6b927464749272))
* fix the values.yaml for news-aggregator chart. Add name for service. ([befd1b7](https://github.com/ivan-poltavskiy/news-aggregator/commit/befd1b7e89487a6a03d5f1e894f17072df735802))

## [2.7.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.6.0...v2.7.0) (2024-10-01)


### Features

* add CloudFormation stack for deploying EKS cluster to AWS. ([a572dd9](https://github.com/ivan-poltavskiy/news-aggregator/commit/a572dd96536161d40fc5e3b6ccdb3839646086a7))

## [2.6.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.5.0...v2.6.0) (2024-10-01)


### Features

* add ability to push the charts and images to ecr. Change Taskfiles structure ([32a9606](https://github.com/ivan-poltavskiy/news-aggregator/commit/32a96066ba870cf41de055a2f81ee4399e9275bb))
* add task for aws auth ([c0ff7c6](https://github.com/ivan-poltavskiy/news-aggregator/commit/c0ff7c6db0788afc62ab4f2d277e6b5d16dd6f23))
* update aws taskfile ([e13194a](https://github.com/ivan-poltavskiy/news-aggregator/commit/e13194a08277817bba64c3039d63818d9d2f842a))
* update aws-action.yml ([22ad5f5](https://github.com/ivan-poltavskiy/news-aggregator/commit/22ad5f5a281b28dcf8484fb797bd7aea4ca9e224))
* update main taskfile. Add docker-build task in push-all-images-to-ecr for building the image before pushing. ([41b9c36](https://github.com/ivan-poltavskiy/news-aggregator/commit/41b9c36e0b194fdb6abd27d6f5fd7e2128be5063))

## [2.5.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.4.0...v2.5.0) (2024-10-01)


### Features

* add webhook for configmap which storage feed groups. ([b59349e](https://github.com/ivan-poltavskiy/news-aggregator/commit/b59349efcc16866a9a816aef00dae90b30d2fdd7))

## [2.4.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.3.0...v2.4.0) (2024-10-01)


### Features

* add cert manager chart for news aggregator server ([b0e034b](https://github.com/ivan-poltavskiy/news-aggregator/commit/b0e034b312bc0f729ecd8cb283915ab42eda29aa))

## [2.3.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.2.0...v2.3.0) (2024-10-01)


### Features

* add GetAllSources endpoint to news aggregator server ([e3ad7b6](https://github.com/ivan-poltavskiy/news-aggregator/commit/e3ad7b6f1b60d6b968e3e36db3d5bfd2042f836e))

## [2.2.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.1.0...v2.2.0) (2024-10-01)


### Features

* add hekm chart for news aggregator server ([d7ccd2f](https://github.com/ivan-poltavskiy/news-aggregator/commit/d7ccd2f8451f4e3b9351377c1d0ec8c3ec850dfc))

## [2.1.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v2.0.0...v2.1.0) (2024-10-01)


### Features

* add taskfile for operator ([fe8578e](https://github.com/ivan-poltavskiy/news-aggregator/commit/fe8578ecf96b1c9f26a73f8b29db650512e94345))

## [2.0.0](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.17...v2.0.0) (2024-10-01)


### ⚠ BREAKING CHANGES

* Add hot news CRD support.

### Features

* Add hot news CRD support. ([f4c948d](https://github.com/ivan-poltavskiy/news-aggregator/commit/f4c948dab79d771e8019d7bfba47ec8386c98de0))


### Bug Fixes

* remove WithEventFilter which filters all watched resources ([0c1e333](https://github.com/ivan-poltavskiy/news-aggregator/commit/0c1e3337a00d528091c6b1577cfb902e5afbf948))

## [1.0.17](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.16...v1.0.17) (2024-08-04)


### Bug Fixes

* fix the tests ([65fc469](https://github.com/ivan-poltavskiy/news-aggregator/commit/65fc469feb0d46befa5260b1ca6d885be56548c0))
* fix the tests ([dd1fdaa](https://github.com/ivan-poltavskiy/news-aggregator/commit/dd1fdaac8d063144d5ae8f88d90acb33bc9a5b10))

## [1.0.16](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.15...v1.0.16) (2024-08-04)


### Bug Fixes

* Update tests, taskfile and actions. Change the mocks generation ([d0a53a7](https://github.com/ivan-poltavskiy/news-aggregator/commit/d0a53a70adc5fe4e9158c9bc42c606e64666f9ed))

## [1.0.15](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.14...v1.0.15) (2024-07-23)


### Bug Fixes

* update github actions configure. ([9d1a116](https://github.com/ivan-poltavskiy/news-aggregator/commit/9d1a11636ea23adaca3a9443f92cf39c12c872a5))

## [1.0.14](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.13...v1.0.14) (2024-07-23)


### Bug Fixes

* delete unnecessary step from release-please.yml ([f50418d](https://github.com/ivan-poltavskiy/news-aggregator/commit/f50418d838a6ec56acfcfc27d2d99cca40d16869))

## [1.0.13](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.12...v1.0.13) (2024-07-22)


### Bug Fixes

* Update release please config ([5dd9531](https://github.com/ivan-poltavskiy/news-aggregator/commit/5dd9531a2376ee85a424ce92fea1e76de669ad2e))

## [1.0.12](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.11...v1.0.12) (2024-07-20)


### Bug Fixes

* Change the name of step. ([a942a73](https://github.com/ivan-poltavskiy/news-aggregator/commit/a942a73959d3807e59548fc6430c95c351da6296))

## [1.0.11](https://github.com/ivan-poltavskiy/news-aggregator/compare/v1.0.10...v1.0.11) (2024-07-19)


### Bug Fixes

* delete debug steps from the release-please.yml ([2e1cb98](https://github.com/ivan-poltavskiy/news-aggregator/commit/2e1cb98f9de32eb3e6867c70b7a3a0a46dc352c1))
