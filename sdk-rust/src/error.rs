use std::{ io, error, fmt };
use std::num::TryFromIntError;
use std::string::FromUtf8Error;

#[derive(Debug)]
pub enum Error {
    FromUtf8(FromUtf8Error),
    TryFromInt(TryFromIntError),
    /// Any [I/O errors](https://doc.rust-lang.org/std/io/struct.Error.html) returned from a Lens function.
    Io(io::Error),
    /// Any [Json errors](https://docs.rs/serde_json/latest/serde_json/struct.Error.html) returned from a Lens function.
    Json(serde_json::Error),
    /// Any [Lens errors](enum.LensError.html) returned from a Lens function.
    Lens(LensError),
}

impl error::Error for Error { }

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            Error::FromUtf8(ref error) => error.fmt(f),
            Error::TryFromInt(ref error) => error.fmt(f),
            Error::Io(ref error) => error.fmt(f),
            Error::Json(ref error) => error.fmt(f),
            Error::Lens(ref error) => error.fmt(f),
        }
    }
}

impl From<FromUtf8Error> for Error {
    fn from(err: FromUtf8Error) -> Error {
        Error::FromUtf8(err)
    }
}

impl From<TryFromIntError> for Error {
    fn from(err: TryFromIntError) -> Error {
        Error::TryFromInt(err)
    }
}

impl From<io::Error> for Error {
    fn from(err: io::Error) -> Error {
        Error::Io(err)
    }
}

impl From<serde_json::Error> for Error {
    fn from(err: serde_json::Error) -> Error {
        Error::Json(err)
    }
}

impl From<LensError> for Error {
    fn from(err: LensError) -> Error {
        Error::Lens(err)
    }
}

#[derive(Clone, PartialEq, Eq, PartialOrd, Ord, Debug, Hash)]
pub enum LensError {
    InputErrorUnsupportedError,
    FailedToWriteErrorToMemError,
}

impl error::Error for LensError { }

impl fmt::Display for LensError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match &*self {
            LensError::InputErrorUnsupportedError => f.write_str("Using errors as inputs is currently unsupported."),
            LensError::FailedToWriteErrorToMemError => f.write_str("An error occured when attempting to write an error to memory."),
        }
    }
}
