This code is basically a rewrite of Martin Kirsche's Wired Logic Sandbox, with some functionality added.

## Command line arguments:
`-scale` <number> - sets the scale of the window (default: 12)

`-width` <number> - sets the width of the window (default: 64)

`-height` <number> - sets the height of the window (default: 64)

`-speed` <number> - sets the speed of the simulation (default: 15)

`-help` - displays the help message

In addition to the command line arguments, the program can load a gif file as a command line argument, and it will be loaded into the program. The gif file must be the last command line argument.

Assuming you have compiled the program into a file called `wired-logic`,here is a sample command line:

`./wired-logic -scale 12 -width 64 -height 64 -speed 15 my-gif.gif`

