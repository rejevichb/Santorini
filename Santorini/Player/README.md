# Player
Contains different player implementations, organized by validity

## Client
Code for Player implementations

## Broken/InfPlace/InfTurn/Valid
`main.go` within each of these subfolders simply allows dynamic loading of the Player creation method, giving the component that loads the plugins the ability to create Players of each type respectively (broken, infinite placement, infinite turn, valid/working as "intended")

## Strategy
Code for Strategy implementations (mapped to the above: broken (sends an invalid turn), infplace (never sends a placement), infturn (never sends a turn), or valid (sends a valid placement and turn))

# Plugin subdirectories
Each of `Valid`, `Broken`, `InfTurn`, and `InfPlace` is a plugin that exports a Player creation function for each of the implementations (rule-abiding, rule-breaking, never providing a turn, and never providing a place respectively)
