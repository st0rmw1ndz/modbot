# modbot

modbot is a seriously over-engineered system for querying different information about your system. Often used with a status bar like [dzen](https://github.com/robm/dzen) or [lemonbar](https://github.com/LemonBoy/bar).

Each part of the output is a module, and modbot is an agregator for these modules. You can query different information about your system like load average, CPU temperature, wireless SSID, and even the output of an arbitrary command or file.

The modules that are in the output are determined at compile time, within the `config.go` file. Don't worry-this process uses a straightforward syntax that doesn't require any programming knowledge. Though, if you're using this, I assume you're smart enough to know that.

modbot is made to be forked, edited, and added onto. If you don't like how something is formatted or looks, you should change it! Go is a very simple and straightforward language, so it should've be too difficult for non-experienced users.

## Attributions

This project was inspired by [gocaudices](https://github.com/LordRusk/gocaudices), which I've used in the past. I created this to be window manager and status bar agnostic, in addition to allowing more options for where the input comes from.

Some ideas for how special cases should be handled were taken from [mblocks](https://gitlab.com/mhdy/mblocks), which is another great project in the same vein as this.
