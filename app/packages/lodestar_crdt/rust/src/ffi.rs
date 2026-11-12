#![allow(clippy::missing_safety_doc)]

use std::slice;

use crate::engine::{BAD_INPUT, Engine};

pub const OK: i32 = 0;

pub struct Handle {
    engine: Engine,
    out: Vec<u8>,
}

unsafe fn bytes<'a>(ptr: *const u8, len: u32) -> &'a [u8] {
    if len == 0 {
        &[]
    } else {
        unsafe { slice::from_raw_parts(ptr, len as usize) }
    }
}

fn code(result: Result<(), i32>) -> i32 {
    result.err().unwrap_or(OK)
}

#[unsafe(no_mangle)]
pub extern "C" fn lc_alloc(len: u32) -> *mut u8 {
    let mut buf = Vec::<u8>::with_capacity(len as usize);
    let ptr = buf.as_mut_ptr();
    std::mem::forget(buf);
    ptr
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_dealloc(ptr: *mut u8, len: u32) {
    if !ptr.is_null() {
        drop(unsafe { Vec::from_raw_parts(ptr, 0, len as usize) });
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn lc_doc_new(client_id: u32) -> *mut Handle {
    Box::into_raw(Box::new(Handle {
        engine: Engine::new(client_id),
        out: Vec::new(),
    }))
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_doc_free(handle: *mut Handle) {
    if !handle.is_null() {
        drop(unsafe { Box::from_raw(handle) });
    }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_out_ptr(handle: *const Handle) -> *const u8 {
    unsafe { (*handle).out.as_ptr() }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_apply(handle: *mut Handle, ptr: *const u8, len: u32) -> i32 {
    let h = unsafe { &*handle };
    code(h.engine.apply(unsafe { bytes(ptr, len) }))
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_state_vector(handle: *mut Handle) -> u32 {
    let h = unsafe { &mut *handle };
    h.out = h.engine.state_vector();
    h.out.len() as u32
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_diff(handle: *mut Handle, ptr: *const u8, len: u32) -> i32 {
    let h = unsafe { &mut *handle };
    match h.engine.diff(unsafe { bytes(ptr, len) }) {
        Ok(out) => {
            h.out = out;
            h.out.len() as i32
        }
        Err(code) => -code,
    }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_text(handle: *mut Handle) -> u32 {
    let h = unsafe { &mut *handle };
    h.out = h.engine.text().into_bytes();
    h.out.len() as u32
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_length(handle: *const Handle) -> u32 {
    unsafe { (*handle).engine.length() }
}

#[unsafe(no_mangle)]
pub unsafe extern "C" fn lc_edit(
    handle: *mut Handle,
    index: u32,
    remove: u32,
    ptr: *const u8,
    len: u32,
) -> i32 {
    let h = unsafe { &*handle };
    let Ok(insert) = std::str::from_utf8(unsafe { bytes(ptr, len) }) else {
        return BAD_INPUT;
    };
    code(h.engine.edit(index, remove, insert))
}
