#!/bin/bash

set -eou pipefail

export YC_TOKEN=$(yc iam create-token --profile=dimasik)
terraform "$@" -var-file=env.tfvars
