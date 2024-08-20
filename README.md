# modbot

modbot is a seriously over-engineered program for querying different information about your system. Often used with a status bar like [dzen](https://github.com/robm/dzen), [lemonbar](https://github.com/LemonBoy/bar), or even the [dwm](https://dwm.suckless.org/) status bar via the `-x` flag.

Each part of the output is a module, and modbot is an agregator for these modules. You can query different information about your system like load average, CPU temperature, wireless SSID, and even the output of an arbitrary command or file.

The modules that are in the output are determined at compile time, within the `config.go` file.

## Attributions

This project was inspired by [gocaudices](https://github.com/LordRusk/gocaudices), which I've used in the past. I created this to be window manager and status bar agnostic, in addition to allowing more options for where the input comes from.

Some ideas for how special cases should be handled were taken from [mblocks](https://gitlab.com/mhdy/mblocks), which is another great project in the same vein as this.

## Usage

```
Usage of modbot:
  -x    set x root window name
```
