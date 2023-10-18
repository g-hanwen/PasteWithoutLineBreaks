## Motivation
Sometimes you copy from a PDF file, and when you paste it into
another file, you get a bunch of additional line breaks that 
you must manually remove. Some PDF readers do this for you, 
but others (e.g. Adobe Acrobat Reader which I use) do not. 

I know that there are so many ways to do this automatically 
(AutoHotKey on Windows, Keyboard Maestro on macOS, Bash script 
written by yourself on Linux, etc.), but if you are lazy, or 
you don't want extra external dependencies, or some of your 
non-technical friends ask you to help them with this, 
then this script is for you. 

## How it works

It substitutes those single occurrences of line breaks and the 
surrounding whitespace with a single space by default. And when 
you are dealing with some languages that do not have a space 
between words (e.g. Chinese), there is also an option to remove 
the line breaks completely. Meanwhile, it preserves continuous 
line breaks, which are usually used to separate paragraphs.

It runs without a Form, so you can just double-click it to run 
it. Toggle the options by right-clicking the tray icon. (For 
those who always get confused by what the menu items mean like 
me, the displayed text is what it will do if you click it 
rather than what it is currently doing.)

## Build
Check the platform-specific notes for [the](https://github.com/golang-design/clipboard)
[dependencies](https://github.com/getlantern/systray). 

The systray package needs CGO, and I don't have a macOS machine, 
so please build it yourself if you are using macOS.

If you want the app icon, you need to install [goversioninfo](https://github.com/josephspurrier/goversioninfo). 
If not so, remove the first line in `main.go`.