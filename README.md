# LensVM

## Building

Once you have the listed prerequisites installed, you should be able to build everything in the repository and run all the tests by running `make test` from the repository root.

### Prerequisites

The following tools need to be installed and added to your PATH before you can build the full contents of the repository:

- Cargo/rustc, typically installed via [rustup](https://www.rust-lang.org/tools/install)
    - Please pay attention to any prerequisites, for example if on Ubuntu you may need to install the `build-essential` package
- The rust wasm32-unknown-unknown compiler, installed using `rustup target add wasm32-unknown-unknown`
- If connection errors are experienced when retrieving rust package dependencies from crates.io, you might need to tweak your `.gitconfig` as per this [comment](https://github.com/rust-lang/cargo/issues/3381#issuecomment-1193730972).
- `npm`, typically installed via [nvm](https://github.com/nvm-sh/nvm#install--update-script)
- [Go](https://golang.google.cn/doc/install)
