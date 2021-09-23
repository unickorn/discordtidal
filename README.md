# discordtidal

Remember when Discord added a Spotify integration and all of your friends started having fun with it, but then being the
weirdo you are, you had TIDAL instead of Spotify Premium?

Well, I do. And this is what I did to make the currently playing TIDAL song show up in my Discord status.

## Okay, so, what is this again?

A tool that makes the TIDAL song you're listening to show up in your Discord status. Only works on Windows using the 
TIDAL desktop application (not the browser version).

Compared to the unsafe branch, this version looks a bit more boring. A picture is below.

![image](https://i.imgur.com/kLXIqSa.png)

## How do I use it?

Unlike the unsafe branch, you can simply download the latest release from the releases tab and run the executable.

While the application window is open, your Discord playing status will be updated to match the song you're playing!

## How does it work?

Unlike the unsafe branch, this version is much simpler, though still cancerous. Here's how it works step by step:

- get the song name from the TIDAL application window title
- look up for the song in TIDAL api
- get the metadata (most importantly the album name)
- connect to discord using [rich-go](https://github.com/hugolgst/rich-go)