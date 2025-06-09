// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use std::error::Error;
use serde::{Serialize, Deserialize};
use lens_sdk::StreamOption;

lens_sdk::define_alloc!();
lens_sdk::define_next!();
lens_sdk::define_transform!(try_transform);

#[derive(Serialize, Deserialize)]
pub struct Value {
    #[serde(rename = "FullName")]
    pub name: String,
    #[serde(rename = "Age")]
	pub age: i64,
}

fn try_transform(
    iter: &mut dyn Iterator<Item = lens_sdk::Result<Option<Value>>>,
) -> Result<StreamOption<Value>, Box<dyn Error>> {
    for item in iter {
        let input = match item? {
            Some(v) => v,
            None => return Ok(StreamOption::None),
        };

        let result = Value {
            name: input.name,
            age: input.age + 1,
        };

        return Ok(StreamOption::Some(result))
    }

    Ok(StreamOption::EndOfStream)
}

#[unsafe(no_mangle)]
pub extern "C" fn inverse() -> *mut u8 {
    match try_inverse() {
        Ok(o) => match o {
            StreamOption::Some(result_json) => lens_sdk::to_mem(lens_sdk::JSON_TYPE_ID, &result_json),
            StreamOption::None => lens_sdk::nil_ptr(),
            StreamOption::EndOfStream => lens_sdk::to_mem(lens_sdk::EOS_TYPE_ID, &[]),
        },
        Err(e) => lens_sdk::to_mem(lens_sdk::ERROR_TYPE_ID, &e.to_string().as_bytes())
    }
}

fn try_inverse() -> Result<StreamOption<Vec<u8>>, Box<dyn Error>> {
    let ptr = unsafe { next() };
    let input = match unsafe { lens_sdk::try_from_mem::<Value>(ptr)? } {
        StreamOption::Some(v) => v,
        // Implementations of `transform` are free to handle nil however they like. In this
        // implementation we chose to return nil given a nil input.
        StreamOption::None => return Ok(StreamOption::None),
        StreamOption::EndOfStream => return Ok(StreamOption::EndOfStream),
    };

    let result = Value {
        name: input.name,
        age: input.age - 1,
    };

    let result_json = serde_json::to_vec(&result)?;
    Ok(StreamOption::Some(result_json))
}
