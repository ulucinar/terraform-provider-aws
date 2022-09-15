#!/usr/bin/env bash

# HO.

LOG_DIR=/Users/alper/data/workspaces/xp/fix-ezgis
(($(dirname "${TERRAFORM_NATIVE_PROVIDER_PATH}")/terraform-provider-aws | tee -a ${LOG_DIR}/native-provider-stdout.log) 3>&1 1>&2 2>&3 | tee -a ${LOG_DIR}/native-provider-stderr.log) 3>&1 1>&2 2>&3
