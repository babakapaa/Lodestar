import 'package:material_ui/material_ui.dart';

import 'contact.dart';

class ContactDraft {
  ContactDraft({
    required this.kind,
    required String value,
    required this.public,
  }) : value = TextEditingController(text: value);

  ContactDraft.from(Contact contact)
    : this(kind: contact.kind, value: contact.value, public: contact.public);

  String kind;
  final TextEditingController value;
  bool public;

  Contact toContact() =>
      Contact(kind: kind, value: value.text.trim(), public: public);

  void dispose() => value.dispose();
}
