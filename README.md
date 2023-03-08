# Untitled Side-Scrolling Game

![untitled-sidescroller-early-demo--03-06](https://user-images.githubusercontent.com/69212809/223170546-fca04225-034e-4587-96b5-c31302b2422c.gif)

A [playable demo](https://shannondybvig.com/project.php?name=sidescroller) is available.


## Status
Currently there are two playable levels, accessible from the world map. Players can collect treasures (for points) and the portal gem (to activate the portal), jump over hazards and creatures (which cause death on contact), and complete the level by exiting through the portal. There is a brief animation when the player dies before they are returned to the world map (if they have lives remaining) or shown a "Game Over" screen.

All artwork are rough stand-ins (I'm learning how to make pixel art alongside building the game).

## Minimum Demo
- [x] Modes for Load, Main Menu, World Map, and Playable Level
- [x] Main Menu
- [x] Save Files

**Playable Level**
- [x] Player Movement (L/R, Jump)
- [x] Collision Logic (rough)
    - [x] Portal (level completion)
    - [x] Brick
    - [x] Quest Item (collect, activate portal)
    - [x] Treasure (collect, +score)
    - [x] Hazard (damage player, level failure)
    - [x] Creature (damage player)
- [ ] Generalize Collision Logic (address collisions between any two objects)
- [ ] Creature Behaviors
    - [x] Movement Logic
    - [ ] Line of Sight
    - [ ] Attack 
- [x] Single-layer background art
- [x] Sprite Sheets (rough)
    - [x] Player Character (does not include jump/fall frames)
    - [x] Brick
    - [x] Quest Item
    - [x] Treasure
    - [x] Hazard
    - [x] Creature
- [x] Treasure Collected Display (either always on screen, or available on [space])
- [ ] Fix Gravity (currently player character floats once jump velocity is resolved)
- [x] Player death animation (rough)
- [ ] Story transitions on levels (entering, dying, exiting)

**World Map**
- [x] Background planet, city placeholders
- [x] Basic movement
- [ ] Simple Menu
    - [x] Enter level from map (by walking character on top of it and pressing a key)
    - [ ] Exit to Main Menu

**Main Menu**
- [x] Start New Game
- [x] Load Game
- [ ] Settings
- [ ] Acknowledgements/Credits
- [x] Exit

## Expanding the Game
- [ ] Pixel-perfect collision
- [ ] Consistent art/aesthetic
- [ ] Story
- [ ] Sound Effects
- [ ] Music
- [ ] Multiple varieties of treasures, hazards, enviro blocks, creatures

## Other Thoughts
I haven't decided yet whether to include things like weapons or jump modifiers (looking at you, pogo stick from the Commander Keen games). I like the idea of building a game where you have a weapon, but very limited ammo (across the entire game) and it is encouraged that you find non-lethal solutions where you can, acknowledging that sometimes it's the only option.

## Influences
This game is heavily influenced by my love of the Commander Keen series (I played both trilogies as a kid and discovered the lost episode in college), Cosmo's Cosmic Adventure, and Cave Story.
