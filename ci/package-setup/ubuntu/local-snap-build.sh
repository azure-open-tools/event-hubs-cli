#!/usr/bin/env bash

# create file snap file based on the template

# --use-lxd (recommended uses inside vms)
snapcraft clean --use-lxd
snapcraft --use-lxd

# submit to snap store