# krishproxy
Simple Websocket Proxy for a Legacy MUD

## Proxying a Legacy MUD

In today’s world, it makes sense to allow playing a MUD directly from a web browser. But when you are faced with a legacy MUD that hasn’t been updated since the turn of the millenium, it can be difficult to update it sufficiently to support websockets.

Instead, this project is a simple proxy service written in Go which will handle the conversion between websockets and the MUD’s innate telnet capabilities.

The code was a quick hack written in about 45 minutes and is just barely sufficient to enable logging into a legacy MUD from a web browser. It “supports” the telnet echo on and echo off sequences and thanks to the xterm.js terminal emulation, ANSI colors and positioning are supported.

If your game does anything beyond those basics, such as using MXP, GMCP or other MUD protocols, you’ll need to extend it to handle the particular sequences that you need and do any additional display work in the fashion most appropriate to your MUD.

I’d appreciate any pull requests or bug reports if you do use this code.

## Try it out

This is currently being used to proxy Krishnak mud.

[Try it out!](http://krishnak.org/game/)
