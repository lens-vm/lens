use std::error::Error;
use serde::{Serialize, Deserialize};

#[derive(Deserialize)]
pub struct Input {
    #[serde(rename(deserialize = "Name"))]
    pub name: String,
    #[serde(rename(deserialize = "Age"))]
	pub age: u64,
}

#[derive(Serialize)]
pub struct Output {
    #[serde(rename(serialize = "FullName"))]
    pub full_name: String,
    #[serde(rename(serialize = "Age"))]
	pub age: u64,
}

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    match try_transform(ptr) {
        Ok(o) => match o {
            Some(result_json) => lens_sdk::to_mem(lens_sdk::JSON_TYPE_ID, &result_json),
            None => lens_sdk::nil_ptr(),
        },
        Err(e) => lens_sdk::to_mem(lens_sdk::ERROR_TYPE_ID, &e.to_string().as_bytes())
    }
}

fn try_transform(ptr: *mut u8) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
    let input = match lens_sdk::try_from_mem::<Input>(ptr)? {
        Some(v) => v,
        // Implementations of `transform` are free to handle nil however they like. In this
        // implementation we chose to return nil given a nil input.
        None => return Ok(None),
    };
    
    let result = Output {
        full_name: input.name,
        age: input.age,
    };
    
    let result_json = serde_json::to_vec(&result)?;
    Ok(Some(result_json))
}
