import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'lti_platform.dart';
import 'lti_platform_dialog.dart';
import 'lti_providers.dart';
import 'tool_urls_section.dart';

class LmsPlatformsScreen extends ConsumerWidget {
  const LmsPlatformsScreen({super.key});

  Future<void> _edit(
    BuildContext context,
    WidgetRef ref, [
    LtiPlatform? platform,
  ]) async {
    final input = await LtiPlatformDialog.show(context, platform);
    if (input == null || !context.mounted) {
      return;
    }
    final api = ref.read(ltiApiProvider);
    await TeamAction.run(
      context,
      () =>
          platform == null ? api.create(input) : api.update(platform.id, input),
    );
    ref.invalidate(ltiPlatformsProvider);
  }

  Future<void> _delete(
    BuildContext context,
    WidgetRef ref,
    LtiPlatform platform,
  ) async {
    final l10n = context.l10n;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.lmsDelete),
        content: Text(l10n.lmsDeleteBody(platform.input.name)),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          FilledButton(
            key: const ValueKey('lms-delete-confirm'),
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.lmsDelete),
          ),
        ],
      ),
    );
    if (confirmed != true || !context.mounted) {
      return;
    }
    await TeamAction.run(
      context,
      () => ref.read(ltiApiProvider).delete(platform.id),
    );
    ref.invalidate(ltiPlatformsProvider);
  }

  List<Widget> _body(
    BuildContext context,
    WidgetRef ref,
    List<LtiPlatform> platforms,
  ) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final tool = ref.watch(ltiToolProvider).value;
    return [
      Text(l10n.lmsHint, style: muted),
      const SizedBox(height: Spacing.l),
      if (tool != null) ...[
        Section(
          title: l10n.lmsToolTitle,
          child: ToolUrlsSection(tool: tool),
        ),
        const SizedBox(height: Spacing.xl),
      ],
      Section(
        title: l10n.lmsPlatforms,
        child: platforms.isEmpty
            ? Text(l10n.lmsNone)
            : Column(
                children: [
                  for (final platform in platforms)
                    LodestarCard(
                      key: ValueKey('lms-${platform.id}'),
                      child: ListTile(
                        contentPadding: EdgeInsets.zero,
                        leading: const Icon(Icons.school_outlined),
                        title: Text(platform.input.name),
                        subtitle: Text(platform.input.issuer),
                        onTap: () => _edit(context, ref, platform),
                        trailing: IconButton(
                          tooltip: l10n.lmsDelete,
                          icon: const Icon(Icons.delete_outline),
                          onPressed: () => _delete(context, ref, platform),
                        ),
                      ),
                    ),
                ],
              ),
      ),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.lmsTitle)),
      floatingActionButton: FloatingActionButton.extended(
        key: const ValueKey('lms-new'),
        onPressed: () => _edit(context, ref),
        icon: const Icon(Icons.add),
        label: Text(l10n.lmsNew),
      ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(ltiPlatformsProvider),
          onRefresh: () => ref.refresh(ltiPlatformsProvider.future),
          builder: (platforms) => _body(context, ref, platforms),
        ),
      ),
    );
  }
}
