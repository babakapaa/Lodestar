use yrs::updates::decoder::Decode;
use yrs::updates::encoder::Encode;
use yrs::{
    Doc, GetString, OffsetKind, Options, ReadTxn, StateVector, Text, TextRef, Transact, Update,
};

pub const BAD_INPUT: i32 = 1;
pub const BAD_RANGE: i32 = 2;

pub struct Engine {
    doc: Doc,
    text: TextRef,
}

impl Engine {
    pub fn new(client_id: u32) -> Self {
        let options = Options {
            client_id: yrs::block::ClientID::new(u64::from(client_id)),
            offset_kind: OffsetKind::Utf16,
            ..Options::default()
        };
        let doc = Doc::with_options(options);
        let text = doc.get_or_insert_text("body");
        Engine { doc, text }
    }

    pub fn apply(&self, update: &[u8]) -> Result<(), i32> {
        let update = Update::decode_v1(update).map_err(|_| BAD_INPUT)?;
        self.doc
            .transact_mut()
            .apply_update(update)
            .map_err(|_| BAD_INPUT)
    }

    pub fn state_vector(&self) -> Vec<u8> {
        self.doc.transact().state_vector().encode_v1()
    }

    pub fn diff(&self, remote: &[u8]) -> Result<Vec<u8>, i32> {
        let remote = StateVector::decode_v1(remote).map_err(|_| BAD_INPUT)?;
        Ok(self.doc.transact().encode_state_as_update_v1(&remote))
    }

    pub fn text(&self) -> String {
        self.text.get_string(&self.doc.transact())
    }

    pub fn length(&self) -> u32 {
        self.text.len(&self.doc.transact())
    }

    pub fn edit(&self, index: u32, remove: u32, insert: &str) -> Result<(), i32> {
        let mut txn = self.doc.transact_mut();
        let length = self.text.len(&txn);
        if index > length || remove > length - index {
            return Err(BAD_RANGE);
        }
        if remove > 0 {
            self.text.remove_range(&mut txn, index, remove);
        }
        if !insert.is_empty() {
            self.text.insert(&mut txn, index, insert);
        }
        Ok(())
    }
}
