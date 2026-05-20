This code is basically a rewrite of [Martin Kirsche's Wired Logic Sandbox](https://github.com/martinkirsche/wired-logic/tree/master/apps/sandbox) using version 2 of the [Ebitengine game engine](https://github.com/hajimehoshi/ebiten), with some functionality added.

## Command line flags and arguments:
### Flags:
`-scale` <number> - sets the scale of the window (default: 12)

`-width` <number> - sets the width of the window (default: 64) 

`-height` <number> - sets the height of the window (default: 64)

`-speed` <number> - sets the speed of the simulation (default: 15)

Note that if you load a gif file (as described in the next section), the `width` and `height` flags will be ignored, and the width and height of the window will be set to the width and height of the gif file.

### Argument:
`<gif file>` - the gif file to load into the program. Must be the last command line argument.

In addition to the command line flags, the program can load a gif file as a command line argument, and it will be loaded into the program. The gif file must be the last command line argument.

Assuming you have compiled the program into a file called `wired-logic`,here is a sample command line:

`./wired-logic -scale 12 -width 64 -height 64 -speed 15`

## Hotkeys:
`Esc` - close the game.

`P` - pause/unpause the simulation. When paused, the simulation will be redrawn in a powered-down state.

`F` - save the current simulation state as a gif file. The file will be saved in the current directory with the name `test.gif`.

`Space` - toggle the pixel under the cursor.

Arrow keys/WASD - move the cursor.