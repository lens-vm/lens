#[derive(Clone, PartialEq, Debug)]
pub enum StreamOption<T> {
    Some(T),
    None,
    EndOfStream,
}

impl <T> StreamOption<T> {
    pub fn ok_or<E>(self, err: E) -> std::result::Result<T, E> {
        match self {
            StreamOption::Some(v) => Ok(v),
            StreamOption::None => Err(err),
            StreamOption::EndOfStream => Err(err),
        }
    }
}
