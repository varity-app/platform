# Varity Scraping Infrastructure
[![Build](https://github.com/VarityPlatform/varity-scraping/actions/workflows/main-workflow.yml/badge.svg)](https://github.com/VarityPlatform/varity-scraping/actions/workflows/main-workflow.yml)

## Cutting a new release
```
# Either on the master branch or in a feature branch run:
git tag v<sem_ver>

# Push the tag
git push origin $(git describe --tags --abbrev=0)
```

## Python Housekeeping
The Makefile contains many helper commands for maintaining python code quality.  Some of the important commands are:
* `make lint` runs the `black` python linter and updates files inside the `src` directory.
* `make pylint` runs the `pylint` linter as a second quality sweep.

## Docker
Various dockerfiles are stored inside the `res` directory.  Images are hosted on [GCR](https://cloud.google.com/container-registry).  To build and publish the images use the `build.sh` script.  Example:
```
# Authenticate with GCR
gcloud auth configure-docker

# Build all images
./build.sh

# Build and publish all images
./build.sh -p

# Build a specific image
./build.sh reddit-scraper
```
More information on building and deploying with GCR can be found [here](https://cloud.google.com/container-registry/docs/quickstart?hl=en_US).

## Deployment
Almost every part of deployment is automated with [Terraform](https://terraform.io).  See more information in the [tf](./tf/README.md) directory.