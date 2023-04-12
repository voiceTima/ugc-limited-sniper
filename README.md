# UGC Limited Sniper Bot

zOMG!!! Another UGC limited sniper bot! 

Features:

- Support for multiple accounts.
- Support for multiple items.
- Is made by a really cool guy.

## Using the project

If you just want to use the bot to snipe UGC, read the [Configuration](#config) section. If you want to don't want to use our prebuilt binaries, you can build the project from source. See the [Building](#building) section for instructions on how to do that.


## Configuration
There is a `config.toml` file in the base directory that should look like the structure below.

```toml
[[roblox.accounts]]
cookie = "your_cookie_here"

[item_ids]
ids = [123456, 789012]
```

If you want to add more accounts, simply add them like so.

```toml
[[roblox.accounts]]
cookie = "your_cookie_here"

[[roblox.accounts]]
cookie = "your_cookie2_here"

[[roblox.accounts]]
cookie = "your_cookie3_here"

[item_ids]
ids = [123456, 789012]
```

You can also add additional item IDs by adding them to the ids array and separating them between a comma.

**Don't forget to save the file before running the program.**

## Downloading 

You can get the compiled binary from the [releases page](https://github.com/PiratePeep/ugc-limited-sniper/releases/tag/windows-release).
## Ratelimits

The Roblox API currently limits API requests to somewhere around `~ 20 requests / min / account`. This means adding more accounts will add longer intervals between checks in order to avoid hitting the ratelimit. This still needs more testing.

This can be avoided by adding proxy support, however, the bot does not have this feature at the moment and I currently have no plans to implement it. If you'd like to see this implemented sooner, please submit a pull request.

## Building

Requirements:
- Go 1.20 or later (could work with older versions, but it hasn't been tested)

### Building for Linux (amd64)

```
GOOS=linux GOARCH=amd64 go build -o ugc_sniper_amd64
```

### Building for Windows (amd64)
```
set GOOS=windows
set GOARCH=amd64
go build -o ugc_sniper_amd64.exe
```

### Building for macOS (amd64)

```
GOOS=darwin GOARCH=amd64 go build -o ugc_sniper_darwin_amd64
```

You can probably figure out the rest if you got this far, nerd. 🤓
