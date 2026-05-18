# rblxutils
roblox modding software. does not fiddle with the client process.

## supported platforms
* windows
sorry sober ppl

## compile and install
rblxutils is a work in progress. the main branch might be unstable.

to hide the console window, helper console window and use the gio backend (usually what you want):
``go install -tags nucular_gio -ldflags="-H windowsgui -X 'main.hide_helper=true'" github.com/axellse/rblxutils@latest``

to show the console window, helper console window and use the gio backend (for debugging):
``go install -tags nucular_gio github.com/axellse/rblxutil@latest``

you can then type `rblxutils` in your terminal to launch the configurator and setup rblxutils. 