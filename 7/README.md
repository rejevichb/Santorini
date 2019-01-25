# Santorini Test Harness

Project 7's deliverable is the `xrules` test harness meant to handle request verification against a given initial boardstate.
The interface for the test harness and valid requests is listed here: [6](http://www.ccs.neu.edu/home/matthias/4500-f18/6.html#%28tech._board._specification%29) and [7](http://www.ccs.neu.edu/home/matthias/4500-f18/7.html).

The implementation of Santorini's board and logic can be found in `Santorini`.

`rules-tests/` contains the `.json` files representing sample player inputs (`n-in.json`) and outputs (`n-out.json`) that can be fed into `xrules` to simulate player actions during a game, and the respective outputs.

To run tests against our `xrules` implementation:
```./TESTME```
