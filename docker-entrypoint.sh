#!/bin/bash
_PORT="${PORT:-80}"
_REPOLOCATION="${REPOLOCATION:-80}"
_IHMLOCATION="${IHMLOCATION:-builds/IHM/ihm.exe}"
_STARTUPLOCATION="${STARTUPLOCATION:-builds/sbRIO-9651/home/lvuser/natinst/bin/startup.rtexe}"
/usr/local/bin/getlaserfile --port=${_PORT} --repolocation=${_REPOLOCATION} --ihmlocation=${_IHMLOCATION} --startuplocation=${_STARTUPLOCATION}