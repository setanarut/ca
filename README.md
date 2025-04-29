# Cellular Automata library in Go

## Cyclic CA Rules

This table presents a collection of cellular automaton rules
along with their creators and rule notations. Each entry follows
the format `Range/Threshold/States/Neighborhood`, where:

- **Range** - defines the distance at which neighboring cells can influence the center cell,
- **Threshold** - specifies the required number of active neighbors for state transitions,
- **States** - indicates the number of possible cell states, and
- **Neighborhood** - (**M** = Moore, **N** = von Neumann) determines the geometric pattern used for neighbor calculations.

| Name             | Author          | Rule       |
| ---------------- | --------------- | ---------- |
| 313              | David Griffeath | `1/3/3/M`  |
| Amoeba           | Jason Rampe     | `3/10/2/N` |
| Black vs White   | Jason Rampe     | `5/23/2/N` |
| Black vs White 2 | Jason Rampe     | `2/5/2/N`  |
| Boiling          | Jason Rampe     | `2/2/6/N`  |
| Bootstrap        | David Griffeath | `2/11/3/M` |
| CCA              | David Griffeath | `1/1/14/N` |
| Cubism           | Jason Rampe     | `2/5/3/N`  |
| Cyclic Spirals   | David Griffeath | `3/5/8/M`  |
| Diamond Spirals  | Jason Rampe     | `1/1/15/N` |
| Fossil Debris    | David Griffeath | `2/9/4/M`  |
| Fuzz             | Jason Rampe     | `4/4/7/N`  |
| Lava Lamp        | Jason Rampe     | `3/15/3/M` |
| Lava Lamp 2      | Jason Rampe     | `2/10/3/M` |
| Maps             | Mirek Wojtowicz | `2/3/5/N`  |
| Perfect Spirals  | David Griffeath | `1/3/4/M`  |
| Stripes          | Mirek Wojtowicz | `3/4/5/N`  |
| Turbulent Phase  | David Griffeath | `2/5/8/M`  |