import 'dart:convert';
import 'dart:typed_data';

class NoteChanges {
  const NoteChanges({
    required this.seq,
    required this.snapshot,
    required this.updates,
  });

  factory NoteChanges.fromJson(Map<String, Object?> json) => NoteChanges(
    seq: (json['seq']! as num).toInt(),
    snapshot: switch (json['snapshot']) {
      final Map<String, Object?> snapshot => base64Decode(
        snapshot['data']! as String,
      ),
      _ => null,
    },
    updates: [
      for (final item in json['updates']! as List<Object?>)
        base64Decode((item! as Map<String, Object?>)['data']! as String),
    ],
  );

  final int seq;
  final Uint8List? snapshot;
  final List<Uint8List> updates;

  List<Uint8List> get all => [?snapshot, ...updates];
}
