use std::mem;
use std::mem::ManuallyDrop;
use std::io::Cursor;
use std::convert::TryInto;
use std::collections::HashMap;
use std::sync::RwLock;
use byteorder::{ReadBytesExt, WriteBytesExt, LittleEndian};

static SRC_PARAM: RwLock<Option<String>> = RwLock::new(None);
static DST_PARAM: RwLock<Option<String>> = RwLock::new(None);

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf);
    return ptr;
}

#[no_mangle]
pub extern fn set_param(id: u32, ptr: *mut u8) {
    let input_str = from_transport_vec(ptr);
    let input_str = serde_json::from_str::<String>(&input_str).unwrap();

    let dst_lock = match id {
        0 => &SRC_PARAM,
        1 => &DST_PARAM,
        _ => panic!("gdfhsand")
    };

    let mut dst = dst_lock.write().unwrap();
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
            return to_transport_vec(&message.as_bytes()).as_mut_ptr()
        },
    };

    let src_param = SRC_PARAM.read().unwrap().clone().unwrap();
    let value = input.get_mut(&src_param).unwrap().clone();
    
    input.remove(&src_param.clone());
    input.insert(DST_PARAM.read().unwrap().clone().unwrap().clone(), value);
    
    let result_json = serde_json::to_vec(&input).unwrap();
    to_transport_vec(&result_json.clone()).as_mut_ptr()
}

fn from_transport_vec(ptr: *mut u8) -> String {
    let len_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr, mem::size_of::<u32>(), mem::size_of::<u32>())
    };

    let rdr = Cursor::new(len_vec);
    let mut rdr = ManuallyDrop::new(rdr);

    let len: usize = rdr.read_u32::<LittleEndian>().unwrap().try_into().unwrap();

    let input_vec: Vec<u8> = unsafe {
        Vec::from_raw_parts(ptr.add(mem::size_of::<u32>()), len, len)
    };
    let input_vec = ManuallyDrop::new(input_vec);
    
    let result = String::from_utf8(input_vec.to_vec()).unwrap().clone();

    mem::drop(input_vec);
    mem::drop(rdr);

    result
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
