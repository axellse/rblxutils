# rblxutils
roblox modding software. does not fiddle with the client process.

## supported platforms
* windows
sorry sober ppl

## compile and install
rblxutils is a work in progress. the main branch might be unstable.

to hide the console window, helper console window and use the gio backend (usually what you want):
``go build -tags nucular_gio,hideHelper -ldflags="-H windowsgui"``

for customized builds, you may use these build tags:
* ``nucular_gio`` - uses the [GIO](https://github.com/gioui/gio) backend. Omitting this tag will use the [shiny](https://pkg.go.dev/golang.org/x/exp/shiny) backend which does not work properly with rblxutils.
* ``hideHelper`` - hides the console window of rblxutils's helper. This flag is only relevant when rblxutils is creating the helper task, so you'll need to remove the helper with task scheduler for any changes to take effect.
* ``keepHelperAlive`` - keeps the helper console window after the helper has exited. Useful for debugging helper panics. This flag is only relevant when rblxutils is creating the helper task, so you'll need to remove the helper with task scheduler for any changes to take effect.

move the binary somewhere, i recommend ``%localappdata%\rblxutils\``. Roblox will be installed in ``./versions/version-xxxxxxxxxxxxxxxx`` relative to the binary's location.

finally run it to launch the configurator. rblxutils will install it's helper and register itself to handle the ``roblox-player`` and ``roblox`` url protocols.                                                    
