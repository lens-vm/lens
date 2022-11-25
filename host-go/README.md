# Lens host-go

Lens/host-go contains a lens host implementation written in Go.

It contains two packages - `engine` is the core lens engine and allows programatic usage of the lens engine.  `config` sits on top of `engine` and allows consumers to provide a lens file containing the configuration of multiple lenses that they wish to be applied to their source data.
