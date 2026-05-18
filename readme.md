# rblxutils
roblox modding software. does not fiddle with the client process.

## supported platforms
* windows
sorry sober ppl

## compile and install
rblxutils is a work in progress. the main branch might be unstable.

first up, clone the repo and cd into it. then build and install:

to hide the console window, helper console window and use the gio backend (usually what you want):
```
go install -tags nucular_gio -ldflags="-H windowsgui -X 'main.hide_helper=true'"
```

to show the console window, helper console window and use the gio backend (for debugging):
```
go install -tags nucular_gio
```

you can then type `rblxutils` in your terminal to launch the configurator and setup rblxutils. 