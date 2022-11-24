use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
pub struct Value {
    #[serde(rename = "FullName")]
    pub name: String,
    #[serde(rename = "Age")]
	pub age: u64,
}

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

const JSON_TYPE_ID: i8 = 1;

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let input = lens_sdk::from_transport_vec::<Value>(ptr).unwrap().unwrap();
    
    let result = Value {
        name: input.name,
        age: input.age + 1,
    };
    
    let result_json = serde_json::to_vec(&result).unwrap();
    lens_sdk::to_transport_vec(JSON_TYPE_ID, &result_json)
}
