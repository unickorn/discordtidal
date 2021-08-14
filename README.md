# discordtidal

### USE AT YOUR OWN RISK! This got my main Discord account banned, as I was carelessly using it thinking Discord would warn me before disabling my account and not responding to the appeal.

### Update: Apparently this wasn't it and the ban was lifted, but still best to be careful!

---

Remember when Discord added a Spotify integration and all of your friends started having fun with it, but then being the
weirdo you are, you had TIDAL instead of Spotify Premium?

Well, I do. And this is what I did to make the currently playing TIDAL song show up in my Discord status.

## Okay, so, what is this again?

A tool that makes the TIDAL song you're listening to show up in your Discord status. Only works on Windows using the 
TIDAL desktop (not the browser version) and looks like this:

![image](https://i.imgur.com/W53wzpq.png)

With a recent Discord update, the background is black no matter the listening/playing status, hence it is even less noticable.

![spotify](https://user-images.githubusercontent.com/29836508/120451801-67df4f00-c39a-11eb-85e7-5173cb47334e.png)
![tidal](https://user-images.githubusercontent.com/29836508/120451811-69a91280-c39a-11eb-8982-5079a4c61412.png)


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
