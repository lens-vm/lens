use std::mem;
use std::mem::ManuallyDrop;
use std::io::{Cursor, Write};
use std::convert::TryInto;
use std::collections::HashMap;
use std::sync::RwLock;
use serde::Deserialize;
use byteorder::{ReadBytesExt, WriteBytesExt, LittleEndian};

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
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf);
    return ptr;
}

#[no_mangle]
pub extern fn set_param(ptr: *mut u8) {
    let input_str = from_transport_vec(ptr);
    let input_str = serde_json::from_str::<Parameters>(&input_str).unwrap();

    let mut dst = PARAMETERS.write().unwrap();
    *dst = Some(input_str);
}

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let input_str = from_transport_vec(ptr);
    let r = serde_json::from_str::<HashMap<String, serde_json::Value>>(&input_str);

    let mut input = match r {
        Ok(i) => i.clone(),
        Err(e) => {
            let message = format!("{}", e);
            return to_transport_vec(ERROR_TYPE_ID, &message.as_bytes()).as_mut_ptr()
        },
    };

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

fn from_transport_vec(ptr: *mut u8) -> String {
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
    
    let result = String::from_utf8(input_vec.to_vec()).unwrap().clone();

    mem::drop(input_vec);
    mem::drop(len_rdr);
    mem::drop(type_rdr);

    result
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
