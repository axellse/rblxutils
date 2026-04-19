# rblxutils
roblox modding software. does not fiddle with the client process.

## state of development?
very early, do not use

## supported platforms
* windows
sorry sober ppl

## compile and install
use the stable branch for stable releases, the main branch is constantly being worked on and may be unstable.

to hide the console/terminal window (usually what you want):
``go build -ldflags -H=windowsgui``

to show it:
``go build``

move the binary somewhere, i recommend ``%localappdata%\rblxutils\``. Roblox will be installed in ``./versions/version-xxxxxxxxxxxxxxxx`` relative to the binary's location.

finally run it to launch the configurator. rblxutils will register itself to handle the ``roblox-player`` and ``roblox`` url protocol.                                                    