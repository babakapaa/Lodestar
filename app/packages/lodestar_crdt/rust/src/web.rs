use wasm_bindgen::prelude::*;

use crate::engine::Engine;

const OK: i32 = 0;

fn code(result: Result<(), i32>) -> i32 {
    result.err().unwrap_or(OK)
}

#[wasm_bindgen]
pub struct CrdtDoc {
    engine: Engine,
}

#[wasm_bindgen]
impl CrdtDoc {
    #[wasm_bindgen(constructor)]
    pub fn new(client_id: u32) -> CrdtDoc {
        CrdtDoc {
            engine: Engine::new(client_id),
        }
    }

    pub fn apply(&self, update: &[u8]) -> i32 {
        code(self.engine.apply(update))
    }

    #[wasm_bindgen(js_name = stateVector)]
    pub fn state_vector(&self) -> Vec<u8> {
        self.engine.state_vector()
    }

    pub fn diff(&self, remote: &[u8]) -> Option<Vec<u8>> {
        self.engine.diff(remote).ok()
    }

    pub fn text(&self) -> String {
        self.engine.text()
    }

    pub fn length(&self) -> u32 {
        self.engine.length()
    }

    pub fn edit(&self, index: u32, remove: u32, insert: &str) -> i32 {
        code(self.engine.edit(index, remove, insert))
    }
}
