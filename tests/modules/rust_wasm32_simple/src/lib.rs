// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use std::error::Error;
use serde::{Serialize, Deserialize};
use lens_sdk::StreamOption;
use lens_sdk::option::StreamOption::{Some, None, EndOfStream};

#[link(wasm_import_module = "lens")]
unsafe extern "C" {
    fn next() -> *mut u8;
}

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

#[unsafe(no_mangle)]
pub extern "C" fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[unsafe(no_mangle)]
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

fn try_transform() -> Result<StreamOption<Vec<u8>>, Box<dyn Error>> {
    let ptr = unsafe { next() };
    let input = match unsafe { lens_sdk::try_from_mem::<Input>(ptr)? } {
        Some(v) => v,
        // Implementations of `transform` are free to handle nil however they like. In this
        // implementation we chose to return nil given a nil input.
        None => return Ok(None),
        EndOfStream => return Ok(EndOfStream)
    };
    
    let result = Output {
        full_name: input.name,
        age: input.age,
    };
    
    let result_json = serde_json::to_vec(&result)?;
    Ok(Some(result_json))
}
