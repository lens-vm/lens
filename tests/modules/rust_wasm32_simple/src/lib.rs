use std::mem;
use std::io::Cursor;
use std::convert::TryInto;
use serde::{Serialize, Deserialize};
use byteorder::{ReadBytesExt, WriteBytesExt, LittleEndian};

#[derive(Deserialize)]
pub struct Input {
    #[serde(rename(deserialize = "Name"))]
    pub name: String,
    #[serde(rename(deserialize = "Age"))]
	pub age: u64,
}

#[derive(Serialize)]
pub struct Result {
    #[serde(rename(serialize = "FullName"))]
    pub full_name: String,
    #[serde(rename(serialize = "Age"))]
	pub age: u64,
}

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf);
    return ptr;
}

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let input_str = from_transport_vec(ptr);
    let r = serde_json::from_str::<Input>(&input_str.clone());
    let input = match r {
        Ok(i) => i,
        Err(e) => {
            let message = format!("{}", e);
            return to_transport_vec(&message.as_bytes()).as_mut_ptr()
        },
    };
    
    let result = Result {
        full_name: input.name,
        age: input.age,
    };
    
    let result_json = serde_json::to_vec(&result).unwrap();
    to_transport_vec(&result_json).as_mut_ptr()
}

fn from_transport_vec(ptr: *mut u8) -> String {
    let len_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr, mem::size_of::<u32>(), mem::size_of::<u32>())
    };
    let mut rdr = Cursor::new(len_vec);
    let len = rdr.read_u32::<LittleEndian>().unwrap().try_into().unwrap();

    mem::forget(rdr);

    let input_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr.add(mem::size_of::<u32>()), len, len)
    };
    
    String::from_utf8(input_vec).unwrap()
}

// we send length-declared strings, not null terminated strings as it makes for safer mem-operations and a slimmer interface
fn to_transport_vec(message: &[u8]) -> Vec::<u8> {
    let len = mem::size_of::<u32>();
    
    let mut wtr = Vec::with_capacity(message.len() + len);
    // cast is coupled to build target! - only slightly though, unless this is a really long string it should be fine.
    // We can also check this before attempting if we want (and/or increase this if needed)
    // LittleEndian is machine specific! (also a pain if users are forced to know this)
    // - actually it doesn't matter, so long as it is consistent with the host's reading
    wtr.write_u32::<LittleEndian>(message.len() as u32).unwrap();
    wtr.extend_from_slice(message);

    wtr
}
