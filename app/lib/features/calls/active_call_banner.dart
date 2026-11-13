import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'calls_providers.dart';

class ActiveCallBanner extends ConsumerWidget {
  const ActiveCallBanner({required this.chatId, super.key});

  final String chatId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final call = ref.watch(activeCallProvider(chatId)).value;
    if (call == null) {
      return const SizedBox.shrink();
    }
    final l10n = context.l10n;
    final palette = context.tokens.palette;
    return Padding(
      padding: const EdgeInsets.fromLTRB(Spacing.m, Spacing.m, Spacing.m, 0),
      child: Material(
        color: palette.accentSoft,
        borderRadius: context.tokens.borderRadius,
        child: ListTile(
          key: const ValueKey('chat-active-call'),
          leading: Icon(Icons.call, color: palette.accent),
          title: Text(l10n.callOngoing),
          subtitle: Text(l10n.callPeople(call.people.length)),
          trailing: FilledButton(
            onPressed: () => context.push(AppRoutes.chatCall(chatId)),
            child: Text(l10n.callJoin),
          ),
        ),
      ),
    );
  }
}
