import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'notes_tabs.dart';

class NotesTabsController extends Notifier<NotesTabs> {
  @override
  NotesTabs build() => const NotesTabs();

  void open(String id, {bool newTab = false}) {
    final tabs = state.open;
    if (tabs.contains(id)) {
      state = NotesTabs(open: tabs, active: id);
      return;
    }
    final active = state.active;
    if (newTab || active == null) {
      state = NotesTabs(open: [...tabs, id], active: id);
      return;
    }
    state = NotesTabs(
      open: [for (final tab in tabs) tab == active ? id : tab],
      active: id,
    );
  }

  void activate(String id) {
    if (state.open.contains(id)) {
      state = NotesTabs(open: state.open, active: id);
    }
  }

  void close(String id) {
    final tabs = state.open;
    final at = tabs.indexOf(id);
    if (at < 0) {
      return;
    }
    final rest = [...tabs]..removeAt(at);
    var active = state.active;
    if (active == id) {
      active = rest.isEmpty ? null : rest[at.clamp(0, rest.length - 1)];
    }
    state = NotesTabs(open: rest, active: active);
  }
}
