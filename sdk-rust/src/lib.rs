/*!
This crate contains members to aid in the construction of a Rust [Lens Module](https://github.com/lens-vm/spec#abi---wasm-module-functions).
*/

use std::mem;
use std::mem::ManuallyDrop;
use std::io::{Cursor, Write};
use serde::Deserialize;
use byteorder::{ReadBytesExt, WriteBytesExt, LittleEndian};

/// [Result](https://doc.rust-lang.org/std/result/enum.Result.html) type alias returned by lens_sdk.
pub mod result;

/// Error types returned by lens_sdk.
pub mod error;

/// Option type returned by lens_sdk.
pub mod option;

/// Alias for an sdk [Error](error/enum.Error.html).
pub type Error = error::Error;

/// Alias for an sdk [Result](result/type.Result.html).
pub type Result<T> = result::Result<T>;

/// Alias for an sdk [StreamOption](option/enum.StreamOption.html).
pub type StreamOption<T> = option::StreamOption<T>;

/// A type id that denotes a simple string-based error.
///
/// If present at the beginning of a byte array being read by a [lens host](https://github.com/lens-vm/lens#Hosts)
/// or [try_from_mem](fn.try_from_mem.html), the byte array will be treated as an error and will be
/// handled accordingly.
pub const ERROR_TYPE_ID: i8 = -1;

/// A type id that denotes a nil value.
pub const NIL_TYPE_ID: i8 = 0;

/// A type id that denotes a json value.
///
/// If present at the beginning of a byte array being read by a [lens host](https://github.com/lens-vm/lens#Hosts)
/// or [try_from_mem](fn.try_from_mem.html), the byte array will be treated as a json value and will be
/// handled accordingly.
pub const JSON_TYPE_ID: i8 = 1;

/// A type id that donates the end of stream.
///
/// If recieved it signals that the end of the stream has been reached and that the source will no longer yield
/// new values.
pub const EOS_TYPE_ID: i8 = i8::MAX;

/// Returns a nil pointer.
///
/// The pointer points to a zeroed byte, and will be interpretted by a [lens host](https://github.com/lens-vm/lens#Hosts)
/// as a nil value.
pub fn nil_ptr() -> *mut u8 {
    &mut 0u8
}

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

/// Read the data held at the given location in memory as the given `TOutput` type.
///
/// The bytes at the given location are expected to be in the correct format for the first (`type_id`) byte.
///
/// `type_id` | Type | Expected format
/// --- | --- | ---
/// < 0 | error | N/A - unsupported, will return an [InputErrorUnsupportedError](error/enum.LensError.html#variant.InputErrorUnsupportedError)
/// 0 | null value | N/A - will return [None](https://doc.rust-lang.org/std/option/enum.Option.html#variant.None)
/// \> 0 | JSON | \[`len`\]\[`json_string`\] where len is the length of bytes in the json_string
///
/// # Safety
///
/// The memory at the given location will be disposed of on read.
///
/// # Errors
///
/// This function will return an [Error](error/enum.Error.html) if the data at the given location is not in the expected
/// format.
pub fn try_from_mem<TOutput: for<'a> Deserialize<'a>>(ptr: *mut u8) -> Result<StreamOption<TOutput>> {
    let type_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr, mem::size_of::<i8>(), mem::size_of::<i8>())
    };

    let type_rdr = Cursor::new(type_vec);
    let mut type_rdr = ManuallyDrop::new(type_rdr);

    let type_id: i8  = type_rdr.read_i8()?;
    if type_id == NIL_TYPE_ID {
        return Ok(StreamOption::None)
    }
    if type_id == EOS_TYPE_ID {
        return Ok(StreamOption::EndOfStream)
    }
    if type_id < 0 {
        return Result::from(error::LensError::InputErrorUnsupportedError)
    }

    let len_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr.add(mem::size_of::<i8>()), mem::size_of::<u32>(), mem::size_of::<u32>())
    };

    let len_rdr = Cursor::new(len_vec);
    let mut len_rdr = ManuallyDrop::new(len_rdr);

    let len: usize = len_rdr.read_u32::<LittleEndian>()?.try_into()?;

    let input_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr.add(mem::size_of::<i8>()+mem::size_of::<u32>()), len, len)
    };
    let input_vec = ManuallyDrop::new(input_vec);

    let json_string = String::from_utf8(input_vec.to_vec())?.clone();

    mem::drop(input_vec);
    mem::drop(len_rdr);
    mem::drop(type_rdr);

    // It is possible for null json values to reach this line, particularly if sourced directly
    // from a 3rd party module, so we ensure that we parse to option as well as the earlier type_id
    // checks.
    Ok(match serde_json::from_str::<Option<TOutput>>(&json_string)? {
        Some(v) => StreamOption::Some(v),
        None => StreamOption::None,
    })
}

/// Write the given `message` bytes to memory, returning a pointer to the first byte.
///
/// Bytes are written in the same format as expected by [lens hosts](https://github.com/lens-vm/lens#Hosts) and
/// [try_from_mem](fn.try_from_mem.html) \- `[type_id][len][json_string]`.
///
/// # Errors
///
/// This function may return the same errors that [io::write](https://doc.rust-lang.org/std/io/trait.Write.html#tymethod.write)
/// may return.
pub fn try_to_mem(type_id: i8, message: &[u8]) -> Result<*mut u8> {
    let buffer = Vec::with_capacity(message.len() + mem::size_of::<u32>() + mem::size_of::<i8>());
    let mut wtr = Cursor::new(buffer);

    wtr.write_i8(type_id)?;
    wtr.set_position(mem::size_of::<i8>() as u64);

    wtr.write_u32::<LittleEndian>(message.len() as u32)?;
    wtr.set_position(mem::size_of::<i8>() as u64 + mem::size_of::<u32>() as u64);
    wtr.write(message)?;

    let result = wtr.into_inner().clone().as_mut_ptr();
    Ok(result)
}

/// Write the given `message` bytes to memory, returning a pointer to the first byte.
///
/// Bytes are written in the same format as expected by [lens hosts](https://github.com/lens-vm/lens#Hosts) and
/// [try_from_mem](fn.try_from_mem.html) \- `[type_id][len][json_string]`.
///
/// This function wraps [try_to_mem](fn.try_to_mem.html), if an error was generated by that internal function call
/// this function will attempt to write the error to memory, if that generates an error, it will attempt to write
/// a [FailedToWriteErrorToMemError](error/enum.LensError.html#variant.FailedToWriteErrorToMemError) to memory,
/// if that fails this function will panic.
///
/// # Panics
///
/// This function will panic if an error was generated internally and all attempts to write it to memory failed.
pub fn to_mem(type_id: i8, message: &[u8]) -> *mut u8 {
    match try_to_mem(type_id, message) {
        Ok(v) => v,
        Err(e) => {
            match try_to_mem(ERROR_TYPE_ID, &e.to_string().as_bytes()) {
                Ok(v) => v,
                Err(_) => {
                    // We have to panic here if this fails, as we can't keep trying and failing to send the error
                    try_to_mem(ERROR_TYPE_ID, error::LensError::FailedToWriteErrorToMemError.to_string().as_bytes()).unwrap()
                }
            }
        }
    }
}
