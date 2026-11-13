import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'link_target.dart';
import 'lti_deep_link.dart';
import 'lti_providers.dart';

class LtiDeepLinkScreen extends ConsumerStatefulWidget {
  const LtiDeepLinkScreen({required this.deepLinkId, super.key});

  final String deepLinkId;

  @override
  ConsumerState<LtiDeepLinkScreen> createState() => _LtiDeepLinkScreenState();
}

class _LtiDeepLinkScreenState extends ConsumerState<LtiDeepLinkScreen> {
  final List<LinkTarget> _chosen = [];
  bool _busy = false;

  bool _isChosen(LinkTarget t) =>
      _chosen.any((c) => c.kind == t.kind && c.id == t.id);

  void _toggle(LtiDeepLink link, LinkTarget t) {
    setState(() {
      if (_isChosen(t)) {
        _chosen.removeWhere((c) => c.kind == t.kind && c.id == t.id);
        return;
      }
      if (!link.acceptMultiple) {
        _chosen.clear();
      }
      _chosen.add(t);
    });
  }

  Future<void> _send() async {
    setState(() => _busy = true);
    await TeamAction.run(context, () async {
      final url = await ref
          .read(ltiApiProvider)
          .answerDeepLink(widget.deepLinkId, _chosen);
      await ref.read(lmsReturnProvider).open(url);
    });
    if (mounted) {
      setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final provider = ltiDeepLinkProvider(widget.deepLinkId);
    return Scaffold(
      appBar: AppBar(title: Text(l10n.lmsDeepLinkTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(provider),
          onRefresh: () => ref.refresh(provider.future),
          builder: (link) => [
            Text(l10n.lmsDeepLinkHint(link.platform), style: muted),
            const SizedBox(height: Spacing.l),
            if (link.targets.isEmpty) Text(l10n.lmsDeepLinkEmpty, style: muted),
            for (final t in link.targets)
              CheckboxListTile(
                key: ValueKey('deep-link-${t.kind}-${t.id}'),
                value: _isChosen(t),
                title: Text(t.title),
                onChanged: _busy ? null : (_) => _toggle(link, t),
              ),
            const SizedBox(height: Spacing.l),
            FilledButton(
              key: const ValueKey('deep-link-send'),
              onPressed: _busy || _chosen.isEmpty ? null : _send,
              child: Text(l10n.lmsDeepLinkSend),
            ),
          ],
        ),
      ),
    );
  }
}
