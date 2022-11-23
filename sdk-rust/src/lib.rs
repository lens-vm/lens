use std::mem;

/// Allocate the given `size` number of bytes in memory and returns a pointer to
/// the first byte.
///
/// The runtime will be instructed to forget this memory, but not dispose of it - the value
/// that is to be held at this location should be written before any other calls are made into
/// the wasm instance.
pub fn alloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf);
    return ptr;
}
