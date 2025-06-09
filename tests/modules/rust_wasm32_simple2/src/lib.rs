// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use std::error::Error;
use serde::{Serialize, Deserialize};
use lens_sdk::StreamOption;

lens_sdk::define_alloc!();
lens_sdk::define_next!();
lens_sdk::define_transform!(try_transform);
lens_sdk::define_inverse!(try_inverse);

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

fn try_inverse(
    iter: &mut dyn Iterator<Item = lens_sdk::Result<Option<Value>>>,
) -> Result<StreamOption<Value>, Box<dyn Error>> {
    for item in iter {
        let input = match item? {
            Some(v) => v,
            None => return Ok(StreamOption::None),
        };

        let result = Value {
            name: input.name,
            age: input.age - 1,
        };

        return Ok(StreamOption::Some(result))
    }

    Ok(StreamOption::EndOfStream)
}
