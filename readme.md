# rblxutils
roblox modding software. does not fiddle with the client process.

## supported platforms
* windows
sorry sober ppl

## compile and install
rblxutils is a work in progress. the main branch might be unstable.

to hide the console/terminal window (usually what you want):
``go build -ldflags -H=windowsgui``

to show it:
``go build``

move the binary somewhere, i recommend ``%localappdata%\rblxutils\``. Roblox will be installed in ``./versions/version-xxxxxxxxxxxxxxxx`` relative to the binary's location.

finally run it to launch the configurator. rblxutils will install it's helper and register itself to handle the ``roblox-player`` and ``roblox`` url protocols.                                                    
