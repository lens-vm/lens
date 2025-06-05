// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use std::collections::HashMap;
use std::sync::RwLock;
use std::error::Error;
use std::{fmt, error};
use serde::Deserialize;
use lens_sdk::StreamOption;
use lens_sdk::option::StreamOption::{Some, None, EndOfStream};

#[link(wasm_import_module = "lens")]
unsafe extern "C" {
    fn next() -> *mut u8;
}

#[derive(Clone, PartialEq, Eq, PartialOrd, Ord, Debug, Hash)]
enum ModuleError {
    ParametersNotSetError,
    PropertyNotFoundError{requested: String},
}

impl error::Error for ModuleError { }

impl fmt::Display for ModuleError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match &*self {
            ModuleError::ParametersNotSetError => f.write_str("Parameters have not been set."),
            ModuleError::PropertyNotFoundError { requested } =>
                write!(f, "The requested property was not found. Requested: {}", requested),
        }
    }
}

#[derive(Deserialize, Clone)]
pub struct Parameters {
    pub src: String,
    pub dst: String,
}

static PARAMETERS: RwLock<StreamOption<Parameters>> = RwLock::new(None);

#[unsafe(no_mangle)]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[unsafe(no_mangle)]
pub extern fn set_param(ptr: *mut u8) -> *mut u8 {
    match try_set_param(ptr) {
        Ok(_) => lens_sdk::nil_ptr(),
        Err(e) => lens_sdk::to_mem(lens_sdk::ERROR_TYPE_ID, &e.to_string().as_bytes())
    }
}

fn try_set_param(ptr: *mut u8) -> Result<(), Box<dyn Error>> {
    let parameter =  unsafe { lens_sdk::try_from_mem::<Parameters>(ptr)? }
        .ok_or(ModuleError::ParametersNotSetError)?;

    let mut dst = PARAMETERS.write()?;
    *dst = Some(parameter);
    Ok(())
}

#[unsafe(no_mangle)]
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
    let ptr = unsafe { next() };
    let mut input = match unsafe { lens_sdk::try_from_mem::<HashMap<String, serde_json::Value>>(ptr)? } {
        Some(v) => v,
        // Implementations of `transform` are free to handle nil however they like. In this
        // implementation we chose to return nil given a nil input.
        None => return Ok(None),
        EndOfStream => return Ok(EndOfStream)
    };

    let params = PARAMETERS.read()?
        .clone()
        .ok_or(ModuleError::ParametersNotSetError)?
        .clone();

    let value = input.get_mut(&params.src)
        .ok_or(ModuleError::PropertyNotFoundError{requested: params.src.clone()})?
        .clone();
    
    input.remove(&params.src);
    input.insert(params.dst, value);
    
    let result_json = serde_json::to_vec(&input)?;
    Ok(Some(result_json))
}
