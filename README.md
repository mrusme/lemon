Lemon
-----

![Lemon D-Bus](lemon-dbus.gif)

[![Static Badge](https://img.shields.io/badge/Donate-Support_this_Project-orange?style=for-the-badge&logo=buymeacoffee&logoColor=%23ffffff&labelColor=%23333&link=https%3A%2F%2Fxn--gckvb8fzb.com%2Fsupport%2F)](https://xn--gckvb8fzb.com/support/) [![Static Badge](https://img.shields.io/badge/Join_on_Matrix-green?style=for-the-badge&logo=element&logoColor=%23ffffff&label=Chat&labelColor=%23333&color=%230DBD8B&link=https%3A%2F%2Fmatrix.to%2F%23%2F%2521PHlbgZTdrhjkCJrfVY%253Amatrix.org)](https://matrix.to/#/%21PHlbgZTdrhjkCJrfVY%3Amatrix.org)

*Lemon* is a *notification proxy* that supports D-Bus, as well as the [Pimoroni 
Raspberry Pi Unicorn Hat 
HD](https://shop.pimoroni.com/products/unicorn-hat-hd?variant=42496126730), 
allowing you to build a cheap, open-source alternative to 
[LaMetric](https://lametric.com).

> **Info:** The `master` branch contains the latest version of *Lemon*. If 
> you're looking for the [legacy 
> version](https://マリウス.com/meet-lemon-a-79-43-open-source-alternative-to-lametric-that-supports-github-ifttt-zapier-webhooks-and-even-integrates-with-pushover/) 
> that was written in Python and was only working on a Raspberry Pi with 
> Pimoroni's Unicorn Hat HD, please refer to [this 
> branch](https://github.com/mrusme/lemon/tree/deprecated).

Lemon can run as a background process (*daemon*) on your Linux/*BSD machine and 
pipe incoming notifications from any supported source to any supported output. 
For example, you can have your *Pushover* notifications displayed as D-Bus 
desktop notifications (*libnotify*).


## Compatibility


### Sources

Lemon currently supports the following notification sources:

- [x] Dummy (for testing purposes)
- [x] [Pushover](https://pushover.net)
- [x] [IFTTT](https://ifttt.com) (via Pushover)


### Outputs 

Lemon currently supports the following outputs:

- [x] D-Bus (desktop notifications)
- [x] Pimoroni Unicorn Hat HD


## Getting started


### From source

```sh
$ go build
$ sudo cp ./lemon /usr/local/bin/lemon
```

### From releases 

Download [the release](https://github.com/mrusme/lemon/releases) for your 
platform, unpack it and move it e.g. to `/usr/local/bin/lemon`.


### Usage

```sh
$ lemon --help
```

#### Pushover

To simplify Pushover setup, this repository contains [a script](login.sh) that 
can be used to obtain the required `Device ID` and `Secret`. The script accepts 
three fixed arguments:

1. Your Pushover e-mail address 
2. Your Pushover password
3. Your Pushover 2FA code

Example:

```sh
./login.sh "$(pass show pushover/email)" "$(pass show pushover/password)" 123456
```


#### Unicorn Hat HD

![Raspberry Pi 3A+ & Pimoroni Unicorn Hat HD](lemon-unicorn.gif)

Download/build the `arm64` version, copy onto the Raspberry Pi and, enable SPI 
via `raspi-config` and use `-outputs unicorn` as output.

