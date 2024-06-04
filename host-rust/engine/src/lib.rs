// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use std::collections::HashMap;
use serde::Deserialize;
use atomic_counter::{AtomicCounter, RelaxedCounter};

#[link(wasm_import_module = "lens")]
extern "C" {
    //module funcs, these can be dynamic as modules are appended? (i don't think that is possible)

    // can still handle those outside of wasm, this can handle passing one item to another - only alloc and transform are req. for that too
    // alloc and transform can be host funcs, not the modules directly - might need module ids/indexes as params
    // this might end up being very slim though, maybe it is not even worth it? Have a look and find out.
/*
    fn alloc(size: usize) -> *mut u8;
    fn transform() -> *mut u8;

    // these are optional - need to figure something out for this! Could fail at wasm init (e.g. runtime.NewInstance)
    fn inverse() -> *mut u8;
    fn set_param(ptr: *mut u8) -> *mut u8;
*/
}

static N_MODULES: Lazy<RelaxedCounter> = Lazy::new(|| RelaxedCounter::new(0));

#[link(wasm_import_module = "lens")]
extern "C" { // to draw next item into pipe from source
    fn next() -> *mut u8;
}

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn append() {
    N_MODULES.inc();
}

#[no_mangle]
pub extern fn transform() -> *mut u8 {// this is wrong and extra
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
    let input = match lens_sdk::try_from_mem::<HashMap<String, serde_json::Value>>(ptr)? {
        Some(v) => v,
        // Implementations of `transform` are free to handle nil however they like. In this
        // implementation we chose to return nil given a nil input.
        None => return Ok(None),
        EndOfStream => return Ok(EndOfStream),
    };


}


////////////////////////////////////////////////////////////////////////


pub extern fn append() {
}

//todo - we don't want the actual runtime in this space, as it is not wasm friendly
pub extern fn new_module(runtime: module::Runtime, path: String) -> Result<module::Module> {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return runtime.NewModule(content)
}

// todo - params type
pub extern fn new_instance(module: module::Module, params: Option<HashMap<String, Any>>) -> Result<module::Instance> {
	return module.NewInstance("transform", paramSets...)
}

// todo - params type
pub extern fn new_inverse(module: module::Module, params: Option<HashMap<String, Any>>) -> Result<module::Instance>  {
	return module.NewInstance("inverse", paramSets...)
}