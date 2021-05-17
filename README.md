# discordtidal

Remember when Discord added a Spotify integration and all of your friends started having fun with it, but then being the
weirdo you are, you had TIDAL instead of Spotify Premium?

Well, I do. And this is what I did to make the currently playing TIDAL song show up in my Discord status.

## Okay, so, what is this again?

A tool that makes the TIDAL song you're listening to show up in your Discord status. Only works on Windows using the 
TIDAL desktop (not the browser version) and looks like this:

![image](https://i.imgur.com/W53wzpq.png)

## How do I use it?

Since this uses a Discord application that's modified to match the song you're listening to, you need to create your own
**Discord Application**.

- Go to https://discord.com/developers, create a new application and set the application ID in the config.toml file.


- Get your Discord account token and set it in the config too.

  You can get this using many methods, I used Firefox' network monitor (Ctrl+Shift+E) to check the https requests that
  contain my token as authorization.


- Get a User Agent of your choice and paste it in the config. I wanted to make this configurable for the fun of it. If
  you don't know what a User Agent is, you can copy your user agent from [here](https://www.whatsmyua.info/).

## How does it work?

That's a little more complicated than you would probably expect, but I had my fun. Here's how it works step by step:

- get the song name from the TIDAL application window title
- look up for the song in TIDAL api
- get the metadata and album cover
- upload the album cover to a Discord application as an application rich presence asset
- update the application name to be the song name (so that it shows up at the top, in bold)
- connect to discord using [rich-go](https://github.com/hugolgst/rich-go)
- store albums and assets in a leveldb so that the same album cover is not uploaded twice (thanks to discord cache)
- clear up assets if they reach 250, since discord has a limit of 300