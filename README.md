A 2D platform action game about a cartoon cat using golang and ebitengine.

Building from source:
Download and unzip source code
If you're on linux, install the ebitengine dependencies as described on their website:
https://ebitengine.org/en/documents/install.html
cd your terminal/powershell to the unzipped source directory
run these:
go mod tidy
go build .
(the executable should be in the source folder)




Controls:
WASD = movement
W = activate door
F = fire flamethrower
SHIFT + A/D = run
SHIFT + W + A/D = long jump
` dev console
Left mouse = edit tile/object

Dev console commands:
fly = fly
tile 1 = enter tile edit mode, and select tile type 1
entity 0 = enter entity edit mode and select type 0, ie. NPCs and enemies

Misc:
fidgets refer to interactive objects like doors and barrels, for lack of a better term

