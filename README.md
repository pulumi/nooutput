# Outputless Pulumi Experiments

A few experiments for alternative programming model projections for working with the Pulumi resource model.

All examples here are in Go, but in principle any of these approaches could be applied to other Pulumi languages.

* [current](./current): The current Pulumi programming model
* [tokenization](./tokenization/): A model where there are no outputs, and "tokens" which have the correct underlying type as the true data as used to smuggle references to eventual values (+ dependencies).
* [blockingget](./blockingget/): A model where accessors block on outputs, and unknowns are embedded into tokens.  This approach doesn't really work, as dependencies are not retained.
* [tokenization-blockingget](./tokenization-blockingget/): A combination of the above, where tokens are used in all cases, but accessors block on outputs being available, so that resolved values of outputs can appear in the token (or `UNKNOWN` for unknowns during preview).  Makes printing and/or stepping through programs substantially "simpler".  In practice, some risk of limiting paralelism by blocking "too soon".