use std::collections::HashMap;
use std::sync::RwLock;
use serde::Deserialize;

#[derive(Deserialize, Clone)]
pub struct Parameters {
    pub src: String,
    pub dst: String,
}

const ERROR_TYPE_ID: i8 = -1;
const JSON_TYPE_ID: i8 = 1;

static PARAMETERS: RwLock<Option<Parameters>> = RwLock::new(None);

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn set_param(ptr: *mut u8) {
    let parameter = lens_sdk::from_transport_vec::<Parameters>(ptr);

    let mut dst = PARAMETERS.write().unwrap();
    *dst = Some(parameter);
}

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let mut input = lens_sdk::from_transport_vec::<HashMap<String, serde_json::Value>>(ptr);

    let params = PARAMETERS.read().unwrap().clone().unwrap();
    let value = match input.get_mut(&params.src) {
        Some(i) => i.clone(),
        None => {
            let message = format!("{} was not found", params.src);
            return lens_sdk::to_transport_vec(ERROR_TYPE_ID, &message.as_bytes())
        },
    };
    
    input.remove(&params.src);
    input.insert(params.dst, value);
    
    let result_json = serde_json::to_vec(&input).unwrap();
    lens_sdk::to_transport_vec(JSON_TYPE_ID, &result_json.clone())
}
