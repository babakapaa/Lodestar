import 'package:material_ui/material_ui.dart';

import 'run_kind.dart';

abstract final class RunIcon {
  static IconData of(RunKind kind) => switch (kind) {
    RunKind.jupyter => Icons.menu_book_outlined,
    RunKind.vscode => Icons.code,
    RunKind.job => Icons.memory,
    RunKind.demo => Icons.slideshow_outlined,
  };
}
