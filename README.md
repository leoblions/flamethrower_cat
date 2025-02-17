# Flamethrower Cat 

## Flamethrower Cat is 2D platform action game about a cartoon cat using golang and ebitengine. It takes inspiration from 80s and 90s platformer games.

## Building from source:
* Download and unzip source code
* Install the go build tools if you haven't. 
![go build tools](https://go.dev/doc/install)
* If you're on linux, install the ebitengine dependencies as described on their website:
![ebitengine installation](https://ebitengine.org/en/documents/install.html)
* cd your terminal/powershell to the unzipped source directory
* run these commands:
* go mod tidy
* go build .
* (the executable should be in the source folder)




## Controls:
* WASD = movement
* W = activate door
* F = fire flamethrower
* SHIFT + A/D = run
* SHIFT + W + A/D = diagonal long jump
* ` dev console
* Left mouse = edit tile/object in edit mode

## Dev console commands:
* fly = fly
* fill = fill tile grid with selected tile 
* tile 1 = enter tile edit mode, and select tile type 1
* entity 0 = enter entity edit mode and select type 0, ie. NPCs and enemies

# Misc Information:

* Fidgets (source code) refer to interactive objects like doors and barrels, for lack of a better term
* Use long jumps with Shift button to get to hard to reach platforms

