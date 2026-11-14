import 'package:flutter/foundation.dart';

@immutable
class QuickSwitchChoice {
  const QuickSwitchChoice.open(String this.noteId) : createPath = null;

  const QuickSwitchChoice.create(String this.createPath) : noteId = null;

  final String? noteId;
  final String? createPath;
}
