mod engine;
#[cfg(not(target_arch = "wasm32"))]
mod ffi;
#[cfg(target_arch = "wasm32")]
mod web;

#[cfg(not(target_arch = "wasm32"))]
pub use ffi::*;

#[cfg(all(test, not(target_arch = "wasm32")))]
mod tests;
