<p align="center">
  <img src="assets/kobowriter.png" />
</p>

# Kobowriter

This small project aims to let you use your old KOBO e-reader (mine is a GLO HD) into a simple, distraction free typewriter.

For years I thought that e-ink was the ultimate medium to write in broad daylight without eye strain of focus fatigue. It seems that others have had the same ideas as we can see in the [Freewrite](https://getfreewrite.com/) or [Pomera](https://www.kickstarter.com/projects/2132003782/pomera-pocket-typewriter-with-e-ink?ref=category_newest&amp;ref=discovery) products.

This project brings the same form factor for a considerably cheaper this way (especially if like me you already have a KOBO at hand).

> Note that the installed software should let you use switch between your normal kobo stock software and the KoboWriter one ; so your kobo is still usable in its default way.

## How it looks like

*Todo*

## How it works

The kobo e-readers have a micro-usb connector to charge and transfer files. With proper kernel modification this usb socket can be used as OTG, letting one plug in any kind of usb device.

Such kernel was compiled by the XCSoar project in order to turn the kobo into a fliying assistant supported by an external GPS. We use their modifications to connect a USB keyboard to the OTG port.

However, the kobo giving no power throught its usb socket, the keyboard has to be powered on its own - you can either use a cheap usb otg power cable [like this one](https://www.amazon.com/AuviPal-Micro-USB-Cable-Power/dp/B07FY9Z9GD/ref=sr_1_3?crid=13TQ5BP3TUJT5&dchild=1&keywords=powered+usb+otg&qid=1630094365&sprefix=powered+%2Caps%2C536&sr=8-3) or modify some keyboard, like I did.

The software lets you use the keyboard to write and edit text files. It's coded in GO, compiled with a toolchain prepared for the KOBO devices, and rely largely on the excellent [FBInk](https://github.com/NiLuJe/FBInk) library to drive the screen, throught its extreemely useful port in Go, [go-fbink](https://github.com/shermp/go-fbink-v2).
## How to build it

> Note that we also provide ready made precompiled binaries for your KOBO

First you need to download and build the **koxtoolchain** on you development computer. This toolchain, once built, will let you build Go program that can run on the KOBO.

## How to install

You can build the software, or we provide pre compiled binaries. You can also use our XCSoar modified installer that will get you the XCSoar program, kernel, and Kobowriter in just one step.