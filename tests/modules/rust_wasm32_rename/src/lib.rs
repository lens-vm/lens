use std::mem;
use std::io::{Cursor, Write};
use std::collections::HashMap;
use std::sync::RwLock;
use serde::Deserialize;
use byteorder::{WriteBytesExt, LittleEndian};

#[derive(Deserialize, Clone)]
pub struct Parameters {
    pub src: String,
    pub dst: String,
}

const ERROR_TYPE_ID: i8 = -1;
const JSON_TYPE_ID: i8 = 1;

static PARAMETERS: RwLock<Option<Parameters>> = RwLock::new(None);

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn set_param(ptr: *mut u8) {
    let parameter = lens_sdk::from_transport_vec::<Parameters>(ptr);

    let mut dst = PARAMETERS.write().unwrap();
    *dst = Some(parameter);
}

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let mut input = lens_sdk::from_transport_vec::<HashMap<String, serde_json::Value>>(ptr);

    let params = PARAMETERS.read().unwrap().clone().unwrap();
    let value = match input.get_mut(&params.src) {
        Some(i) => i.clone(),
        None => {
            let message = format!("{} was not found", params.src);
            return to_transport_vec(ERROR_TYPE_ID, &message.as_bytes()).as_mut_ptr()
        },
    };
    
    input.remove(&params.src);
    input.insert(params.dst, value);
    
    let result_json = serde_json::to_vec(&input).unwrap();
    to_transport_vec(JSON_TYPE_ID, &result_json.clone()).as_mut_ptr()
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
