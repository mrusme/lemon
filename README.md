Lemon
-----

LED Monitor, for monitoring all sort of stuff and getting notified when something happens.

![GitHub Commit](docs/github-commit.gif)

[More examples here.](https://github.com/mrusme/lemon/tree/master/docs)

## Requirements

- WiFi
- A Raspberry Pi with a [Pimoroni Unicorn Hat HD](https://shop.pimoroni.de/products/unicorn-hat-hd)
- Raspbian on its microSD card
- A GitHub account

## Installation

On your Raspberry:

```bash
$ raspi-config nonint do_spi 0
$ reboot
$ aptitude install python3 python3-pip python3-dev python3-spidev libtiff5-dev libjpeg-dev zlib1g-dev libfreetype6-dev liblcms2-dev libwebp-dev libharfbuzz-dev libfribidi-dev tcl8.6-dev tk8.6-dev python-tk
$ pip3 install gunicorn falcon pillow ujson unicornhathd py-pushover-open-client
$ cd /opt
$ git clone https://github.com/mrusme/lemon.git
$ ln -s /opt/lemon/init.d/lemon /etc/init.d/lemon
$ update-rc.d lemon defaults
```

Lemon will now be run automatically every time your Raspberry starts.

## Running manually

```bash
$ cd lemon/
$ ./lemon.sh
```

## Configuration

### Port (on boot)

You can change the port by creating a file named `/etc/lemon` and adding the following content to it:

```bash
export PORT=1337
```

The init.d-script will look for the file and in case it was found source it.

### Port (manually)

```bash
PORT=1337 ./lemon.sh
```

## Integrations

### GitHub Webhooks

First, make Lemon's port (e.g. `20001`) accessible from the internet. You'll probably need to add a NAT rule to your firewall/router or fire up ngrok or something like that.
Then, set up hooks on the GitHub repos (`https://github.com/mrusme/<REPO>/settings/hooks`) or organisations (`https://github.com/organizations/<ORG>/settings/hooks`) that you'd like to receive notifications for. 
For the `Payload URL`, use `http://<YOUR-DNS-NAME>:<FORWARDED-PORT>/github`. As `Content type` use `application/json`. The secret is not necessary (yet). Also, configure the webhook to `Send me everything` and don't forget to check the `Active` checkbox.
That's it. You should be receiving notifications for at least a handful of implemented GitHub events.

### IFTTT

You can create all sort of notifications using [IFTTT's Maker Webhooks](https://ifttt.com/maker_webhooks). Configuration is pretty easy:

Build your IFTTT rule by specifying the `If this` service and using `Webhooks` as `then that` service. Configure the webhook like this:

![Webhook](docs/ifttt-make_a_web_request.png)

### Zapier

Like with IFTTT, you can send webhooks from Zapier to Lemon. Configuration for Zapier is pretty easy as well:

Build your Zapier zap and use `Webhooks` as an action. Configure the webhook like this:

![Webhook POST](docs/zapier-webhook-01.png)

![Webhook Payload](docs/zapier-webhook-02.png)

### Pushover

Lemon provides an integration into Pushover using the Open Client API. Therefor, a separate config needs to be created, which will be read from/written to by Lemon. See [lemon-pushover.cfg](lemon-pushover.cfg) for an example configuration. In order to activate the Pushover plug-in, the environment variable `LEMON_PUSHOVER_CONFIG` needs to be set to the location of your `lemon-pushover.cfg`, e.g. `LEMON_PUSHOVER_CONFIG=/etc/lemon-pushover.cfg`. Make sure that the file is read and writable by the user you run Lemon under!

Lemon will register a new device in your Pushover account (named `lemon`). This device can be targeted by other Pushover clients and it will of course also receive all un-targeted notifications. Be aware that you'll need a [Pushover desktop license](https://pushover.net/clients) in order to use this feature. However, they do provide a 7-day-trial for you to test it.

If don't want to fiddle around with DynDNS, NAT or ngrok in order to make Lemon's HTTP port reachable from GitHub, IFTTT, Zapier and other webhook providers, you can set up Lemon to only use Pushover, which doesn't require you to expose any port. The Pushover client implementation uses a websocket to connect to the Pushover API and retrieve notifications. It basically acts like a web browser, hence you'll be able to use it even within networks you have no/little control over.

**Info regarding 2FA (2 Factor Authentication)**: For Pushover Open Client API integration, Lemon uses [jonogreenz/py-pushover-open-client](https://github.com/jonogreenz/py-pushover-open-client/). However, this library did not natively support 2FA, but a [PR was set](https://github.com/jonogreenz/py-pushover-open-client/pull/6) to the author of the library to include a possibility for an initial 2FA authentication. As soon as this PR was accepted and a new version of the library was released, 2FA will be available here as well.

## API

You can also attach any other service by using the generic Lemon API, which is accessible through `http://yourdns:20001/api`. The body of the requests/webhooks should contain the following JSON:

```json
{ 
  "icon": "icon-name-here",
  "text": "The text to be scrolled through after the icon animation finished"
}
```

You need to make sure that the icon name you specify actually exists within the `icons/` folder and is a `.png` file. If you want to use the icon `icons/docker.png`, simply specify `"icon": "docker"` in the JSON.

By default, each icon animation is being repeated three times. However, you can override that by specifying an additional JSON parameter named `icon_repeat`, e.g.:

```json
{ 
  "icon": "doom-bloody",
  "icon_repeat": 1,
  "text": "Game over!"
}
```

Animation cycle time can also be adjusted individually, in order to allow smoother playback of longer animations:

```json
{ 
  "icon": "microsoft",
  "icon_repeat": 1,
  "icon_cycle_time": "0.01",
  "text": "Micro$oft"
}
```

The default font can be adjusted by specifying the `text_font` property. In order for this to work, the font needs to be available as `.ttf` file inside the `fonts/` directory. Example:

```json
{ 
  "icon": "psy",
  "text": "Oppa Gangnam Style!",
  "text_font": "Hack-Regular"
}
```

## Testing

```bash
curl -X "POST" "http://raspberrypi:20001/ifttt" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "icon": "youtube",
  "text": "New video on YouTube!"
}'
```

## Kudos to ...

- [@pimoroni](https://github.com/pimoroni) for their awesome hardware
- [@source-foundry](https://github.com/source-foundry) for [Hack](https://github.com/source-foundry/Hack), the best Termin/Editor font there has ever been

## "Let me tell you..."

Sure, [tell me](https://twitter.com/intent/tweet?text=@mrusme%20regarding%20Lemon,%20let%20me%20tell%20you%20that...)!
