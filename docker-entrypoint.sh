#!/bin/bash
_PORT="${PORT:-80}"
_CONFIGLOCATION="${CONFIGLOCATION:-/etc/getlaserfile/config.yaml}"
/usr/local/bin/getlaserfile --port=${_PORT} --config=${_CONFIGLOCATION}