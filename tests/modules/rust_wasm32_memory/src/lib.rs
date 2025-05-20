// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use std::error::Error;
use serde::{Deserialize, Serialize};
use lens_sdk::StreamOption;
use lens_sdk::option::StreamOption::{Some, None, EndOfStream};

#[link(wasm_import_module = "lens")]
extern "C" {
    fn next() -> *mut u8;
}

#[no_mangle]
pub extern "C" fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern "C" fn transform() -> *mut u8 {
    match try_transform() {
        Ok(o) => match o {
            Some(result_json) => lens_sdk::to_mem(lens_sdk::JSON_TYPE_ID, &result_json),
            None => lens_sdk::nil_ptr(),
            EndOfStream => lens_sdk::to_mem(lens_sdk::EOS_TYPE_ID, &[]),
        },
        Err(e) => lens_sdk::to_mem(lens_sdk::ERROR_TYPE_ID, &e.to_string().as_bytes())
    }
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Input {
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename= "__type")]
	pub __type: String,
    #[serde(rename = "Array")]
	pub array: Vec<String>,
}

fn try_transform() -> Result<StreamOption<Vec<u8>>, Box<dyn Error>> {
    let ptr = unsafe { next() };
    let mut input = match lens_sdk::try_from_mem::<Input>(ptr)? {
        Some(v) => v,
        // Implementations of `transform` are free to handle nil however they like. In this
        // implementation we chose to return nil given a nil input.
        None => return Ok(None),
        EndOfStream => return Ok(EndOfStream)
    };

    for _ in 1..1000 {
        // Copy the input a bunch of times to make sure that the memory is not cleaned up
        // before we are finished with it.
        input = input.clone();
    }

    let result_json = serde_json::to_vec(&input.clone())?;
    lens_sdk::free_transport_buffer(ptr)?;
    Ok(Some(result_json))
}
