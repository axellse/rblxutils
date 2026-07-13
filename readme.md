<h1 align="center">
  <img src="https://raw.githubusercontent.com/axellse/rblxutils/refs/heads/main/pictures/logo.png" height="100"></img>
  <br>
  rblxutils
</h1>
<p align="center">
    <i>free and open source roblox modloader & bootstrapper</i>
    <br>
    <a href="https://github.com/axellse/ray/releases">latest release</a>
    <span> | </span>
    <a href="https://github.com/axellse/rblxutils/actions/workflows/build.yml">github actions builds</a>
</p>

# introduction
Rblxutils (Roblox -> rblx + utilities -> utils) is a free and open source Roblox modloader and bootstrapper. As a bootstrapper, it serves as one of the few non-Bloxstrap forks, being written in go. The bootstrapper comes with log file reading (activity tracking) which aids other features such as Discord Rich Presence and server history.

<img src="https://raw.githubusercontent.com/axellse/rblxutils/refs/heads/main/pictures/bootstrapper.png" alt="screenshot of bootstrapper" height="300"></img>

The main feature of Rblxutils is the modloader. Rblxutils applies mods using the open [Roblox Community Modding Format (RCMF)](https://rcmf.axell.me/v1/) standard to which it also acts as a reference implementation. Unlike Bloxstrap and other software that can only apply mods to Roblox files on disk, Rblxutils can modify any asset using a proxy around Roblox's [assetdelivery service](https://create.roblox.com/docs/cloud/reference/domains/assetdelivery).

RCMF provides several benefits over traditional Roblox mods. Apart from the aforementioned asset modification capabilities, RCMF separates each mod into its own isolated package, making it very easy to swap, remove or quickly toggle different mods. RCMF also allows easier and more standardized distribution of mods using ``.rcmf`` files.

Rblxutils can also import traditional Roblox file mods as zips, as long as the zip specifies where to place the files using subdirectories (eg. if you can drop the mod into Bloxstrap's mod directory, you can import it into Rblxutils). In this case, Rblxutils will translate the zip into an RCMF which it will then import, giving you all the same isolation benefits.

<img src="https://raw.githubusercontent.com/axellse/rblxutils/refs/heads/main/pictures/configurator.png" alt="screenshot of configurator" height="300"></img>

## getting started
Rblxutils is Windows-only at the moment. I'm planning on adding RoL support via Sober in the future, but this is not something I've looked into yet.

You can download a stable release binary here on Github in the [releases section](https://github.com/axellse/rblxutils/releases) or alternatively a freshly built, possibly unstable binary from [Github actions](https://github.com/axellse/rblxutils/actions/workflows/build.yml).

Rblxutils can run anywhere if it finds a config file (``config.json``) in the current directory, but will launch its installer if that is not the case. Rblxutils's installer is very easy to use and will install to ``%LocalAppData%\rblxutils``, as well as create a start menu shortcut for easy access. Rblxutils auto-updates once it's installed. 

Launching Rblxutils will send you to the configurator where you can install mods and configure settings. When ran the first time Rblxutils will install its helper and register itself to handle the ``roblox-player`` and ``roblox`` url protocols. Roblox will be installed in ``./versions/version-xxxxxxxxxxxxxxxx`` relative to the binary's location.
## compile yourself
You can also compile Rblxutils yourself if you want. Since the main branch is the working branch and might be unstable, you can download a source code zip from the releases tab. If you're doing development or just want the latest and greatest stuff, go ahead and clone the repo like usual.

Rblxutils is written in Go, so you'll need to [download and install that](https://go.dev/doc/install).

For a standard build which is usually what you want and what the release binaries use, use the following command:

```
go build -tags nucular_gio,hideHelper -ldflags="-H windowsgui"
```

For customized builds, you may use these build tags:
* ``nucular_gio`` - uses the [GIO](https://github.com/gioui/gio) backend. Omitting this tag will use the [shiny](https://pkg.go.dev/golang.org/x/exp/shiny) backend which does not work properly with rblxutils.
* ``hideHelper`` - hides the console window of Rblxutils's helper. Rblxutils will modify the task scheduler task if it does not match this configured setting.
* ``keepHelperAlive`` - keeps the helper console window after the helper has exited. Useful for debugging helper panics. Rblxutils will modify the task scheduler task if it does not match this configured setting.

Once you've got your binary, you can follow the [getting started section](#getting-started). Do note though that any updates will override your custom build with a standard one.