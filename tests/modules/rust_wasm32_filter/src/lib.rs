use std::error::Error;
use std::iter::Iterator;
use serde::{Serialize, Deserialize};
use lens_sdk::StreamOption;

#[link(wasm_import_module = "lens")]
unsafe extern "C" {
    fn next() -> *mut u8;
}

#[derive(Serialize, Deserialize)]
#[cfg_attr(test, derive(PartialEq, Debug))]
pub struct Value {
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename = "__type")]
	pub type_name: String,
}

#[unsafe(no_mangle)]
pub extern "C" fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[unsafe(no_mangle)]
pub extern "C" fn transform() -> *mut u8 {
    lens_sdk::next(|| -> *mut u8 { unsafe { next() } }, try_transform)
}

fn try_transform(
    iter: &mut dyn Iterator<Item = lens_sdk::Result<Option<Value>>>,
) -> Result<StreamOption<Value>, Box<dyn Error>> {
    for item in iter {
        let input = match item? {
            Some(v) => v,
            None => continue,
        };

        if input.type_name == "pass" {
            return Ok(StreamOption::Some(input))
        }
    }

    Ok(StreamOption::EndOfStream)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_try_transform_pass() {
        let input = [
            Ok(
                Some(
                    Value{
                        name: "John".to_string(),
                        type_name: "pass".to_string(),
                    },
                ),
            ),
        ];

        let result = try_transform(&mut input.into_iter()).unwrap();

        assert_eq!(
            result,
            StreamOption::Some(Value{
                name: "John".to_string(),
                type_name: "pass".to_string(),
            }),
        );
    }

    #[test]
    fn test_try_transform_skip() {
        let input = [
            Ok(
                Some(
                    Value{
                        name: "Fred".to_string(),
                        type_name: "fail".to_string(),
                    },
                ),
            ),
        ];

        let result = try_transform(&mut input.into_iter()).unwrap();

        assert_eq!(
            result,
            StreamOption::EndOfStream,
        );
    }
}
