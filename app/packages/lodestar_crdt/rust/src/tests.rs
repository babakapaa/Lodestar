use std::slice;

use crate::engine::{BAD_INPUT, BAD_RANGE};
use crate::ffi::*;

struct Doc2(*mut Handle);

impl Doc2 {
    fn new(id: u32) -> Self {
        Doc2(lc_doc_new(id))
    }
    fn out(&self, len: usize) -> Vec<u8> {
        unsafe { slice::from_raw_parts(lc_out_ptr(self.0), len).to_vec() }
    }
    fn edit(&self, index: u32, remove: u32, text: &str) -> i32 {
        unsafe { lc_edit(self.0, index, remove, text.as_ptr(), text.len() as u32) }
    }
    fn text(&self) -> String {
        let len = unsafe { lc_text(self.0) } as usize;
        String::from_utf8(self.out(len)).unwrap()
    }
    fn sv(&self) -> Vec<u8> {
        let len = unsafe { lc_state_vector(self.0) } as usize;
        self.out(len)
    }
    fn diff(&self, sv: &[u8]) -> Vec<u8> {
        let len = unsafe { lc_diff(self.0, sv.as_ptr(), sv.len() as u32) };
        self.out(len as usize)
    }
    fn apply(&self, update: &[u8]) -> i32 {
        unsafe { lc_apply(self.0, update.as_ptr(), update.len() as u32) }
    }
    fn sync_from(&self, other: &Doc2) {
        let update = other.diff(&self.sv());
        assert_eq!(self.apply(&update), OK);
    }
}

impl Drop for Doc2 {
    fn drop(&mut self) {
        unsafe { lc_doc_free(self.0) }
    }
}

#[test]
fn offline_edits_merge() {
    let a = Doc2::new(1);
    let b = Doc2::new(2);
    assert_eq!(a.edit(0, 0, "Клетка состоит из ядра."), OK);
    b.sync_from(&a);
    assert_eq!(a.edit(0, 0, "Айгерим: "), OK);
    let end = unsafe { lc_length(b.0) };
    assert_eq!(b.edit(end, 0, " И мембраны 🧬."), OK);
    a.sync_from(&b);
    b.sync_from(&a);
    assert_eq!(a.text(), b.text());
    assert_eq!(a.text(), "Айгерим: Клетка состоит из ядра. И мембраны 🧬.");
}

#[test]
fn rejects_bad_input() {
    let a = Doc2::new(1);
    assert_eq!(a.edit(5, 0, "x"), BAD_RANGE);
    assert_eq!(a.apply(&[255, 255, 255]), BAD_INPUT);
    assert!(unsafe { lc_diff(a.0, [255u8, 255].as_ptr(), 2) } < 0);
}

#[test]
fn utf16_offsets() {
    let a = Doc2::new(1);
    a.edit(0, 0, "a🧬b");
    assert_eq!(unsafe { lc_length(a.0) }, 4);
    assert_eq!(a.edit(1, 2, ""), OK);
    assert_eq!(a.text(), "ab");
}

#[test]
fn many_peers_converge() {
    let mut seed: u64 = 42;
    let mut next = move |n: u32| {
        seed = seed
            .wrapping_mul(6364136223846793005)
            .wrapping_add(1442695040888963407);
        ((seed >> 33) as u32) % n.max(1)
    };
    for _ in 0..100 {
        let docs: Vec<Doc2> = (1..=3).map(Doc2::new).collect();
        for _ in 0..60 {
            let d = &docs[next(3) as usize];
            let len = unsafe { lc_length(d.0) };
            if len > 0 && next(2) == 0 {
                let at = next(len);
                let n = 1 + next((len - at).min(3));
                assert_eq!(d.edit(at, n.min(len - at), ""), OK);
            } else {
                assert_eq!(d.edit(next(len + 1), 0, "ёж"), OK);
            }
            if next(5) == 0 {
                let (i, j) = (next(3) as usize, next(3) as usize);
                if i != j {
                    docs[j].sync_from(&docs[i]);
                }
            }
        }
        for i in 0..3 {
            for j in 0..3 {
                if i != j {
                    docs[j].sync_from(&docs[i]);
                }
            }
        }
        assert_eq!(docs[0].text(), docs[1].text());
        assert_eq!(docs[1].text(), docs[2].text());
    }
}
