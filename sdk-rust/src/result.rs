use crate::error::{LensError, Error};

pub type Result<T> = std::result::Result<T, Error>;

impl<T> From<LensError> for Result<T> {
    fn from(err: LensError) -> Result<T> {
        Err(Error::Lens(err))
    }
}
