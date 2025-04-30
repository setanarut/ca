# Cellular Automata library in Go

## Cyclic CA Rules

This table presents a collection of cellular automaton rules
along with their creators and rule notations. Each entry follows
the format `Range/Threshold/States/Neighborhood`, where:

- **Range** - defines the distance at which neighboring cells can influence the center cell,
- **Threshold** - specifies the required number of active neighbors for state transitions,
- **States** - indicates the number of possible cell states, and
- **Neighborhood** - (**M** = Moore, **N** = von Neumann) determines the geometric pattern used for neighbor calculations.

| Rule     | Name             | Author          |
| -------- | ---------------- | --------------- |
| 1/3/3/M  | 313              | David Griffeath |
| 3/10/2/N | Amoeba           | Jason Rampe     |
| 5/23/2/N | Black vs White   | Jason Rampe     |
| 2/5/2/N  | Black vs White 2 | Jason Rampe     |
| 2/2/6/N  | Boiling          | Jason Rampe     |
| 2/11/3/M | Bootstrap        | David Griffeath |
| 1/1/14/N | CCA              | David Griffeath |
| 2/5/3/N  | Cubism           | Jason Rampe     |
| 3/5/8/M  | Cyclic Spirals   | David Griffeath |
| 1/1/15/N | Diamond Spirals  | Jason Rampe     |
| 2/9/4/M  | Fossil Debris    | David Griffeath |
| 4/4/7/N  | Fuzz             | Jason Rampe     |
| 3/15/3/M | Lava Lamp        | Jason Rampe     |
| 2/10/3/M | Lava Lamp 2      | Jason Rampe     |
| 2/3/5/N  | Maps             | Mirek Wojtowicz |
| 1/3/4/M  | Perfect Spirals  | David Griffeath |
| 3/4/5/N  | Stripes          | Mirek Wojtowicz |
| 2/5/8/M  | Turbulent Phase  | David Griffeath |

## Life-like CA Rules

| Rule          | Name               | Description                                                                                                                                                                                                                                     |
| :------------ | :----------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| B1357/S1357   | Replicator         | Edward Fredkin's replicating automaton: every pattern is eventually replaced by multiple copies of itself.                                                                                                                                      |
| B2/S          | Seeds              | All patterns are phoenixes, meaning that every live cell immediately dies, and many patterns lead to explosive chaotic growth. However, some engineered patterns with complex behavior are known.                                               |
| B25/S4        |                    | This rule supports a small self-replicating pattern which, when combined with a small glider pattern, causes the glider to bounce back and forth in a pseudorandom walk.                                                                        |
| B3/S012345678 | Life without Death | Also known as Inkspot or Flakes. Cells that become alive never die. It combines chaotic growth with more structured ladder-like patterns that can be used to simulate arbitrary Boolean circuits.                                               |
| B3/S23        | Life               | Highly complex behavior.                                                                                                                                                                                                                        |
| B34/S34       | 34 Life            | Was initially thought to be a stable alternative to Life, until computer simulation found that larger patterns tend to explode. Has many small oscillators and spaceships.                                                                      |
| B35678/S5678  | Diamoeba           | Forms large diamonds with chaotically fluctuating boundaries. First studied by Dean Hickerson, who in 1993 offered a $50 prize to find a pattern that fills space with live cells; the prize was won in 1999 by David Bell.                     |
| B36/S125      | 2x2                | If a pattern is composed of 2x2 blocks, it will continue to evolve in the same form; grouping these blocks into larger powers of two leads to the same behavior, but slower. Has complex oscillators of high periods as well as a small glider. |
| B36/S23       | HighLife           | Similar to Life but with a small self-replicating pattern.                                                                                                                                                                                      |
| B3678/S34678  | Day & Night        | Symmetric under on-off reversal. Has engineered patterns with highly complex behavior.                                                                                                                                                          |
| B368/S245     | Morley             | Named after Stephen Morley; also called Move. Supports very high-period and slow spaceships.                                                                                                                                                    |
| B4678/S35678  | Anneal             | Also called the twisted majority rule. Symmetric under on-off reversal. Approximates the curve-shortening flow on the boundaries between live and dead cells.                                                                                   |