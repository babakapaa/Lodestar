import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/lodestar_card.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'assignment_source.dart';
import 'link_target.dart';
import 'lti_context.dart';
import 'lti_providers.dart';

class LtiContextCard extends ConsumerStatefulWidget {
  const LtiContextCard({
    required this.context,
    required this.onChanged,
    this.openLink = true,
    super.key,
  });

  final LtiContext context;
  final VoidCallback onChanged;
  final bool openLink;

  @override
  ConsumerState<LtiContextCard> createState() => _LtiContextCardState();
}

class _LtiContextCardState extends ConsumerState<LtiContextCard> {
  bool _busy = false;

  Future<void> _run(Future<void> Function() action) async {
    setState(() => _busy = true);
    await TeamAction.run(context, action);
    if (mounted) {
      setState(() => _busy = false);
    }
    widget.onChanged();
  }

  Future<void> _sync() async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    await _run(() async {
      final report = await ref.read(ltiApiProvider).sync(widget.context.id);
      messenger.showSnackBar(
        SnackBar(
          content: Text(
            l10n.lmsSynced(
              report.members,
              report.enrolled,
              report.works,
              report.scores,
            ),
          ),
        ),
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final c = widget.context;
    final link = c.link;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final api = ref.read(ltiApiProvider);
    return LodestarCard(
      key: ValueKey('lti-context-${c.id}'),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(c.name, style: Theme.of(context).textTheme.titleMedium),
          Text(l10n.lmsFrom(c.platform), style: muted),
          const SizedBox(height: Spacing.m),
          if (link == null)
            Text(l10n.lmsNotLinked)
          else
            ListTile(
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.link),
              title: Text(l10n.lmsLinkedTo(link.title)),
              onTap: widget.openLink ? () => context.go(link.path) : null,
              trailing: c.instructor
                  ? TextButton(
                      key: const ValueKey('lti-unlink'),
                      onPressed: _busy
                          ? null
                          : () => _run(() => api.unlink(c.id)),
                      child: Text(l10n.lmsUnlink),
                    )
                  : null,
            ),
          if (c.instructor && link == null) ...[
            const SizedBox(height: Spacing.s),
            if (c.options.isEmpty)
              Text(l10n.lmsNoOptions, style: muted)
            else
              for (final option in c.options)
                ListTile(
                  key: ValueKey('lti-link-${option.id}'),
                  contentPadding: EdgeInsets.zero,
                  leading: const Icon(Icons.add_link),
                  title: Text(option.title),
                  subtitle: Text(l10n.lmsLinkHint),
                  enabled: !_busy,
                  onTap: () => _link(option),
                ),
          ],
          if (link != null) ...[
            const SizedBox(height: Spacing.m),
            Text(l10n.lmsSource, style: Theme.of(context).textTheme.titleSmall),
            const SizedBox(height: Spacing.s),
            if (c.sourceEditable)
              SegmentedButton<AssignmentSource>(
                key: const ValueKey('lti-source'),
                segments: [
                  for (final source in AssignmentSource.values)
                    ButtonSegment(
                      value: source,
                      label: Text(_sourceName(context, source)),
                    ),
                ],
                selected: {c.source},
                onSelectionChanged: _busy
                    ? null
                    : (value) => _run(() => api.source(c.id, value.single)),
              )
            else
              Text(_sourceName(context, c.source)),
            Text(_sourceHint(context, c.source), style: muted),
          ],
          if (c.syncError != null) ...[
            const SizedBox(height: Spacing.m),
            NoticeBanner(
              icon: Icons.sync_problem_outlined,
              message: l10n.lmsSyncFailed,
            ),
          ],
          if (c.canSync) ...[
            const SizedBox(height: Spacing.m),
            Row(
              children: [
                Expanded(
                  child: Text(
                    c.syncedAt == null
                        ? l10n.lmsNeverSynced
                        : l10n.lmsSyncedAt(
                            DateFormat.MMMd(
                              Localizations.localeOf(context).toLanguageTag(),
                            ).add_Hm().format(c.syncedAt!),
                          ),
                    style: muted,
                  ),
                ),
                TextButton.icon(
                  key: const ValueKey('lti-sync'),
                  onPressed: _busy ? null : _sync,
                  icon: const Icon(Icons.sync),
                  label: Text(l10n.lmsSync),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }

  Future<void> _link(LinkTarget option) =>
      _run(() => ref.read(ltiApiProvider).link(widget.context.id, option));

  String _sourceName(BuildContext context, AssignmentSource source) =>
      switch (source) {
        AssignmentSource.lodestar => context.l10n.lmsSourceLodestar,
        AssignmentSource.lms => context.l10n.lmsSourceLms,
      };

  String _sourceHint(BuildContext context, AssignmentSource source) =>
      switch (source) {
        AssignmentSource.lodestar => context.l10n.lmsSourceLodestarHint,
        AssignmentSource.lms => context.l10n.lmsSourceLmsHint,
      };
}
