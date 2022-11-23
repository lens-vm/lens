use serde::{Serialize, Deserialize};

#[derive(Deserialize)]
pub struct Input {
    #[serde(rename(deserialize = "Name"))]
    pub name: String,
    #[serde(rename(deserialize = "Age"))]
	pub age: u64,
}

#[derive(Serialize)]
pub struct Result {
    #[serde(rename(serialize = "FullName"))]
    pub full_name: String,
    #[serde(rename(serialize = "Age"))]
	pub age: u64,
}

const JSON_TYPE_ID: i8 = 1;

#[no_mangle]
pub extern fn alloc(size: usize) -> *mut u8 {
    lens_sdk::alloc(size)
}

#[no_mangle]
pub extern fn transform(ptr: *mut u8) -> *mut u8 {
    let input = lens_sdk::from_transport_vec::<Input>(ptr);
    
    let result = Result {
        full_name: input.name,
        age: input.age,
    };
    
    let result_json = serde_json::to_vec(&result).unwrap();
    lens_sdk::to_transport_vec(JSON_TYPE_ID, &result_json)
}
