use std::mem;
use std::io::{Cursor, Write};
use serde::{Serialize, Deserialize};
use byteorder::{WriteBytesExt, LittleEndian};

#[derive(Serialize, Deserialize)]
pub struct Value {
    #[serde(rename = "FullName")]
    pub name: String,
    #[serde(rename = "Age")]
	pub age: u64,
}

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

const JSON_TYPE_ID: i8 = 1;

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let input = lens_sdk::from_transport_vec::<Value>(ptr);
    
    let result = Value {
        name: input.name,
        age: input.age + 1,
    };
    
    let result_json = serde_json::to_vec(&result).unwrap();
    to_transport_vec(JSON_TYPE_ID, &result_json).as_mut_ptr()
}

// we send length-declared strings, not null terminated strings as it makes for safer mem-operations and a slimmer interface
fn to_transport_vec(type_id: i8, message: &[u8]) -> Vec::<u8> {
    let buffer = Vec::with_capacity(message.len() + mem::size_of::<u32>() + mem::size_of::<i8>());
    let mut wtr = Cursor::new(buffer);

    wtr.write_i8(type_id).unwrap();
    wtr.set_position(mem::size_of::<i8>() as u64);
    // cast is coupled to build target! - only slightly though, unless this is a really long string it should be fine.
    // We can also check this before attempting if we want (and/or increase this if needed)
    // LittleEndian is machine specific! (also a pain if users are forced to know this)
    // - actually it doesn't matter, so long as it is consistent with the host's reading
    wtr.write_u32::<LittleEndian>(message.len() as u32).unwrap();
    wtr.set_position(mem::size_of::<i8>() as u64 + mem::size_of::<u32>() as u64);
    wtr.write(message).unwrap();

    wtr.into_inner().clone()
}
