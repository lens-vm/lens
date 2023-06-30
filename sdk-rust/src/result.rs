// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

use crate::error::{LensError, Error};

pub type Result<T> = std::result::Result<T, Error>;

impl<T> From<LensError> for Result<T> {
    fn from(err: LensError) -> Result<T> {
        Err(Error::Lens(err))
    }
}
