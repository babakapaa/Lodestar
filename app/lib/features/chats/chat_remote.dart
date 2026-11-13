import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../l10n/l10n_context.dart';

abstract final class ChatRemote {
  static Future<bool> run(
    BuildContext context,
    Future<void> Function() request,
  ) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      await request();
      return true;
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
      return false;
    }
  }
}
