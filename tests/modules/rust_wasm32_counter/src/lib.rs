use std::error::Error;
use serde::{Serialize, Deserialize};
use atomic_counter::{AtomicCounter, RelaxedCounter};
use once_cell::sync::Lazy;
use lens_sdk::StreamOption;

#[link(wasm_import_module = "lens")]
unsafe extern "C" {
    fn next() -> *mut u8;
}

fn next_ptr() -> *mut u8 {
    unsafe {
        next()
    }
}

#[derive(Serialize, Deserialize, Default)]
pub struct Value {
    #[serde(rename = "Id")]
	pub id: usize,
    #[serde(rename = "Name")]
    pub name: String,
}

static COUNTER: Lazy<RelaxedCounter> = Lazy::new(|| RelaxedCounter::new(0));

#[unsafe(no_mangle)]
pub extern "C" fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[unsafe(no_mangle)]
pub extern "C" fn transform() -> *mut u8 {
    lens_sdk::next(next_ptr, try_transform)
}

fn try_transform(
    iter: &mut dyn Iterator<Item = lens_sdk::Result<Option<Value>>>,
) -> Result<StreamOption<Value>, Box<dyn Error>> {
    let input = match iter.next() {
        Some(v) => v?,
        None => return Ok(StreamOption::EndOfStream)
    };

    COUNTER.inc();

    let result = Value {
        id: COUNTER.get(),
        name: input.unwrap_or_default().name,
    };

    Ok(StreamOption::Some(result))
}
