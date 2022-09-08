# discordtidal

> ### :warning: USE AT YOUR OWN RISK!
> This code is openly abusing Discord's API to upload album covers as "game assets" in real time to a Discord application, which probably wouldn't be good if Discord caught your account doing that. I take no responsibility if your account gets terminated using discordtidal!

Check out the safer branch [here](https://github.com/Unickorn/discordtidal/tree/safer).
It offers the following:
- less risk of your account getting terminated
- much simpler to use, just click on an exe and it works
- looks much more boring

---

Remember when Discord added a Spotify integration and all of your friends started having fun with it, but then being the
weirdo you are, you had TIDAL instead of Spotify Premium?

Well, I do. And this is what I did to make the currently playing TIDAL song show up in my Discord status.

## Okay, so, what is this again?

A tool that makes the TIDAL song you're listening to show up in your Discord status. Only works on Windows using the 
TIDAL desktop application (not the browser version) and looks like this compared to a Spotify listening status:


![spotify](https://user-images.githubusercontent.com/29836508/120451801-67df4f00-c39a-11eb-85e7-5173cb47334e.png)
![tidal](https://user-images.githubusercontent.com/29836508/120451811-69a91280-c39a-11eb-8982-5079a4c61412.png)


## How do I use it?

Since this uses a Discord application that's modified to match the song you're listening to, you need to create your own
**Discord Application**.

- Download a release from the releases tab.

- Run the exe, it will most likely crash but that's fine. Now you will see a `config.toml` file in the folder.

- Go to https://discord.com/developers, create a new application and set the Application ID value in the config.

- Get your Discord account token and set it in the config too.

  You can get this using different methods, I used Firefox' network monitor (Ctrl+Shift+E) to check the https requests that
  contain my token as authorization.

- Get a User Agent of your choice and paste it in the config. I wanted to make this configurable for the fun of it. If
  you don't know what a User Agent is, you can copy your user agent from [here](https://www.whatsmyua.info/).
  
- Start Discord, then discordtidal, and then the TIDAL desktop app. You may need to add it as a game for it to show up on your playing status.

- Start listening!

> :warning: The album artwork might not appear the first time you're listening to a new album, because Discord's asset cache takes some time to update. In a couple minutes or the next time you listen to the same album, the cover should be there.

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
