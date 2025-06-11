# LensVM

LensVM is a bi-directional transformation engine originally built for [DefraDB](https://github.com/sourcenetwork/defradb) but available as a standalone tool. It enables the transformation of data in both forward and reverse directions, via user-created Lenses written in any language and compiled to WASM.

Each Lens is executed in its own independent WASM environment, allowing for pipelines to be safely constructed from Lenses sourced from multiple parties.

## Using LensVM

To use LensVM, you will need two things, first you will need to chose how you wish to interact with LensVM, via an [engine](#using-a-lensvm-engine).

Second, you will need to source one or more Lenses, containing the WASM [transforms](#writing-lenses) that you wish to apply to your data.

### Using a LensVM engine

There is currently a single Golang LensVM engine available for use.  Implementations (or wrappers) in other languages are planned, if you have a strong interest in using Lens via another language besides Go please thumb-up or create a github issue for it.

The documentation for using the Go engine is [here](host-go/README.md).

### Writing Lenses

Lenses may be written in any language you chose, so long as they are compiled into valid WASM.

Each Lens must contain the following functions:
- `alloc(unsigned64)` - This exported function is mandatory, and will be called by the LensVM engine to allocate a memory block of the given size.
- `next() unsigned8` - This imported function is mandatory, and will allow the Lens to pull in a pointer the next data-item from the LensVM engine.
- `set_param(unsigned8) unsigned8` - This exported function is optional, it can be provided if you wish to provide static configured data on engine initialization to the Lens.  It receives a pointer to the configured data, and returns a pointer to an ok/error response. It will be called once, before the Lens receives any data items from the engine.
- `transform() unsigned8` - This exported function is mandatory, it is the function in which the data will be transformed.  It can pull zero-many items from `next()`, transform the inputs, then return a single pointer to the transformed output item.  Only one output can be yielded at a time, but the Lens can be stateful if desired, allowing you to cache multiple outputs and yield them one by one.
- `inverse() unsigned8` - This exported function is optional, and allows you to define the inverse of `transform()` should you wish - it is otherwise defined in exactly the same way as `transform()`.

Data is sent across the WASM boundary (to `set_param()`, `transform()`, `inverse()`) using the following format:
```
[TypeId][Length][Payload]
```
`TypeId` is a signed 8-byte integer. `Length` is an unsigned 32 byte integer.  `TypeId` can have the following values:
- `-1`, this indicates an error. A `Length` may be set, and error string may be provided as a `Payload`.
- `0`, this indicates a nil item. `Length` and `Payload` will not exist.
- `1`, this indicates a json item. `Length` will be set to the length of the `Payload`, and `Payload` will contain the json serialized item.
- `127`, this indicates the end of the data stream, and that there are no more items to pull from `next()`, or return from `transform()` or `inverse()`. `Length` and `Payload` will not exist.

There is a Rust SDK to aid the writing of Lenses in Rust, it should eliminate the need to worry about all of the above. It is can be found on [GitHub](/sdk-rust) and [crates.io](https://crates.io/crates/lens_sdk).

There are example Lenses written in AssemblyScript and Rust in [this repo](/tests/modules), and [DefraDB](https://github.com/sourcenetwork/defradb/tree/develop/tests/lenses).

## Building

Once you have the listed prerequisites installed, you should be able to build everything in the repository and run all the tests by running `make test` from the repository root.

### Prerequisites

The following tools need to be installed and added to your PATH before you can build the full contents of the repository:

- [rustup](https://www.rust-lang.org/tools/install) and Cargo/rustc, typically installed via rustup.
    - Please pay attention to any prerequisites, for example if on Ubuntu you may need to install the `build-essential` package
- If connection errors are experienced when retrieving rust package dependencies from crates.io, you might need to tweak your `.gitconfig` as per this [comment](https://github.com/rust-lang/cargo/issues/3381#issuecomment-1193730972).
- `npm`, typically installed via [nvm](https://github.com/nvm-sh/nvm#install--update-script)
- [Go](https://golang.google.cn/doc/install)
