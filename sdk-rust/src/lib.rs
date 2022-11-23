use std::mem;
use std::mem::ManuallyDrop;
use std::io::{Cursor, Write};
use serde::Deserialize;
use byteorder::{ReadBytesExt, WriteBytesExt, LittleEndian};

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
/// The bytes at the given location are expected to be in the format `[type_id][len][json_string]`.
///
/// # Safety
///
/// The memory at the given location will be disposed of on read.
pub fn from_transport_vec<TOutput: for<'a> Deserialize<'a>>(ptr: *mut u8) -> TOutput {
    let type_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr, mem::size_of::<i8>(), mem::size_of::<i8>())
    };

    let type_rdr = Cursor::new(type_vec);
    let mut type_rdr = ManuallyDrop::new(type_rdr);

    let type_id: i8  = type_rdr.read_i8().unwrap().try_into().unwrap();
    if type_id <= 0 {
        panic!("todo - we only support type > 0 atm, return error")
    }

    let len_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr.add(mem::size_of::<i8>()), mem::size_of::<u32>(), mem::size_of::<u32>())
    };

    let len_rdr = Cursor::new(len_vec);
    let mut len_rdr = ManuallyDrop::new(len_rdr);

    let len: usize = len_rdr.read_u32::<LittleEndian>().unwrap().try_into().unwrap();

    let input_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr.add(mem::size_of::<i8>()+mem::size_of::<u32>()), len, len)
    };
    let input_vec = ManuallyDrop::new(input_vec);

    let json_string = String::from_utf8(input_vec.to_vec()).unwrap().clone();

    mem::drop(input_vec);
    mem::drop(len_rdr);
    mem::drop(type_rdr);

    serde_json::from_str::<TOutput>(&json_string).unwrap()
}

/// Write the given `message` bytes to memory, returning a pointer to the first byte.
///
/// Bytes are written in the same format as expected by [lens hosts](https://github.com/lens-vm/lens#Hosts) and
/// [from_transport_vec](fn.from_transport_vec.html) \- `[type_id][len][json_string]`.
pub fn to_transport_vec(type_id: i8, message: &[u8]) -> *mut u8 {
    let buffer = Vec::with_capacity(message.len() + mem::size_of::<u32>() + mem::size_of::<i8>());
    let mut wtr = Cursor::new(buffer);

    wtr.write_i8(type_id).unwrap();
    wtr.set_position(mem::size_of::<i8>() as u64);

    wtr.write_u32::<LittleEndian>(message.len() as u32).unwrap();
    wtr.set_position(mem::size_of::<i8>() as u64 + mem::size_of::<u32>() as u64);
    wtr.write(message).unwrap();

    wtr.into_inner().clone().as_mut_ptr()
}
