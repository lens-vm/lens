use std::error::Error;
use std::collections::VecDeque;
use std::sync::RwLock;
use serde::{Serialize, Deserialize};
use once_cell::sync::Lazy;
use lens_sdk::StreamOption;
use lens_sdk::option::StreamOption::{Some, None, EndOfStream};

#[link(wasm_import_module = "lens")]
extern "C" {
    fn next() -> *mut u8;
}

#[derive(Serialize, Deserialize)]
pub struct Book {
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename = "PageNumbers")]
	pub page_numbers: Vec<i32>,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Page {
    #[serde(rename = "BookName")]
    pub book_name: String,
    #[serde(rename = "Number")]
	pub number: i32,
}

static PENDING_PAGES: RwLock<Lazy<VecDeque<Page>>> = RwLock::new(Lazy::new(|| VecDeque::new()));

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn transform() -> *mut u8 {
    match try_transform() {
        Ok(o) => match o {
            Some(result_json) => lens_sdk::to_mem(lens_sdk::JSON_TYPE_ID, &result_json.clone()),
            None => lens_sdk::nil_ptr(),
            EndOfStream => lens_sdk::to_mem(lens_sdk::EOS_TYPE_ID, &[]),
        },
        Err(e) => lens_sdk::to_mem(lens_sdk::ERROR_TYPE_ID, &e.to_string().as_bytes())
    }
}

fn try_transform() -> Result<StreamOption<Vec<u8>>, Box<dyn Error>> {
    let mut pending_pages = PENDING_PAGES.write()?;

    if pending_pages.len() == 0 {
        let ptr = unsafe { next() };
        let input = match lens_sdk::try_from_mem::<Book>(ptr)? {
            Some(v) => v,
            // Implementations of `transform` are free to handle nil however they like. In this
            // implementation we chose to return nil given a nil input.
            None => return Ok(None),
            EndOfStream => return Ok(EndOfStream),
        };

        for page_number in input.page_numbers {
            let page = Page {
                book_name: input.name.clone(),
                number: page_number,
            };
            pending_pages.push_back(page);
        }
        lens_sdk::free_transport_buffer(ptr)?;
    }

    if pending_pages.len() > 0 {
        let result = pending_pages.pop_front();
        let result_json = serde_json::to_vec(&result)?;
        return Ok(Some(result_json));
    }

    Ok(None)
}
