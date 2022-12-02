use std::error::Error;
use serde::{Serialize, Deserialize};
use lens_sdk::StreamOption;
use lens_sdk::option::StreamOption::{Some, None, EndOfStream};

#[link(wasm_import_module = "lens")]
extern "C" {
    fn next() -> *mut u8;
}

#[derive(Serialize, Deserialize)]
pub struct Value {
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename = "__type")]
	pub type_name: String,
}

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn transform() -> *mut u8 {
    match try_transform() {
        Ok(o) => match o {
            Some(result_json) => lens_sdk::to_mem(lens_sdk::JSON_TYPE_ID, &result_json),
            None => lens_sdk::nil_ptr(),
            EndOfStream => lens_sdk::to_mem(lens_sdk::EOS_TYPE_ID, &[]),
        },
        Err(e) => lens_sdk::to_mem(lens_sdk::ERROR_TYPE_ID, &e.to_string().as_bytes())
    }
}

fn try_transform() -> Result<StreamOption<Vec<u8>>, Box<dyn Error>> {
    loop {
        let ptr = unsafe { next() };
        let input = match lens_sdk::try_from_mem::<Value>(ptr)? {
            Some(v) => v,
            // Implementations of `transform` are free to handle nil however they like. In this
            // implementation we chose to return nil given a nil input.
            None => return Ok(None),
            EndOfStream => return Ok(EndOfStream),
        };

        if input.type_name == "pass" {
            let result_json = serde_json::to_vec(&input)?;
            return Ok(Some(result_json))
        }
        lens_sdk::free_transport_buffer(ptr)?;
    }
}
