# Side-Scrolling Pixel Game

## Status
Currently there is a rough semi-playable level (that cannot be won, but which can cause instant death), a world map (also rough), and a load screen. All artwork are rough stand-ins (I'm learning how to make pixel art alongside building the game).

## Minimum Demo
- [x] Modes for Load, Main Menu, World Map, and Playable Level
- [x] Main Menu
- [ ] Save Files

**Playable Level**
- [x] Player Movement (L/R, Jump)
- [ ] Collision Logic (rough)
    - [x] Portal (level completion)
    - [x] Brick
    - [x] Quest Item (collect, activate portal)
    - [x] Treasure (collect, +score)
    - [x] Hazard (damage player, level failure)
    - [ ] Creature (damage player)
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
- [ ] Treasure Collected Display (either always on screen, or available on [space])
- [ ] Fix Gravity (currently player character floats once jump velocity is resolved)
- [ ] Option to Exit Level without completing

**World Map**
- [x] Background planet, city placeholders
- [x] Basic movement
- [ ] Simple Menu
    - [x] Enter level from map (by walking character on top of it and pressing a key)
    - [ ] Exit to Main Menu

**Main Menu**
- [x] Start New Game
- [ ] Load Game
- [ ] Settings
- [ ] Acknowledgements/Credits
- [x] Exit

## Expanding the Game
Once a basic structure is in place, I can think about the following:
- Creating multiple types of creatures, treasures, building blocks, hazards, etc.
- Developing an actual story and consistent aesthetic (and so creating entirely new art)
- Sounds effects and music

## Other Thoughts
I haven't decided yet whether to include things like weapons or jump modifiers (looking at you, pogo stick from the Keen games). I like the idea of building a game where you have a weapon, but very limited ammo (across the entire game) and it is encouraged that you find non-lethal solutions where you can, acknowledging that sometimes it's the only option.

## Influences
This effort is heavily influenced by my love of the Commander Keen series (as a kid I played both trilogies, and I discovered the lost episode in college), Cosmo's Cosmic Adventure, and Cave Story.
