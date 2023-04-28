# wmtgr

Check [webmention.io][] for new mentions and send updates to Telegram.

[webmention.io]: https://webmention.io

## Usage
```
NAME:
   wmtgr - webmentions to telegram relay

USAGE:
   wmtgr [global options] command [command options] [arguments...]

DESCRIPTION:
   wmtgr periodically checks webmention.io for new webmentions and sends them to Telegram using Telegram bot API.

AUTHOR:
   Tymoteusz Makowski <tymek.makowski@gmail.com>

COMMANDS:
   fetch    fetch webmentions once, print, and exit
   relay    fetch webmentions periodically and relay them to Telegram
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --webmention-token value  webmention.io API token
   --help, -h                show help
```

### Relay
This CLI comes with a command to continuously check for webmentions and send messages via Telegram bot API.
```
NAME:
   wmtgr relay - fetch webmentions periodically and relay them to Telegram

USAGE:
   wmtgr relay [command options] [arguments...]

OPTIONS:
   --telegram-token value      Telegram bot API token
   --telegram-chat-id value    Telegram chat ID
   --interval value, -n value  interval between checks (default: 1h0m0s)
   --help, -h                  show help
```

## Name
wmtgr - \[w\]eb\[m\]entions \[t\]ele\[g\]ram \[r\]elay
