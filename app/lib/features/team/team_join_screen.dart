import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'team_card.dart';
import 'team_labels_context.dart';
import 'team_providers.dart';

class TeamJoinScreen extends ConsumerStatefulWidget {
  const TeamJoinScreen({required this.code, super.key});

  final String code;

  @override
  ConsumerState<TeamJoinScreen> createState() => _TeamJoinScreenState();
}

class _TeamJoinScreenState extends ConsumerState<TeamJoinScreen> {
  String? _error;
  bool _busy = false;

  Future<void> _join() async {
    final l10n = context.l10n;
    final router = GoRouter.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final details = await ref.read(teamApiProvider).join(widget.code);
      ref.invalidate(myTeamsProvider);
      unawaited(router.pushReplacement(AppRoutes.team(details.team.id)));
    } on ApiException catch (error) {
      setState(() => _error = ApiErrorText.of(l10n, error));
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  List<Widget> _body(TeamCard card) {
    final l10n = context.l10n;
    final textTheme = Theme.of(context).textTheme;
    final error = _error;
    return [
      LodestarCard(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(card.team.name, style: textTheme.titleLarge),
            const SizedBox(height: Spacing.s),
            Text([card.scope.context, card.scope.title].join('\n')),
            const SizedBox(height: Spacing.s),
            Text(context.teamLabels.rules(card.scope)),
            const SizedBox(height: Spacing.s),
            Text(card.team.names),
          ],
        ),
      ),
      if (error != null) FormError(message: error),
      const SizedBox(height: Spacing.xl),
      LodestarButton(
        key: const ValueKey('team-join'),
        label: l10n.teamJoin,
        icon: Icons.group_add_outlined,
        onPressed: _busy ? null : _join,
      ),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final preview = ref.watch(teamPreviewProvider(widget.code));
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.teamJoinTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: preview,
          onRefresh: () => ref.refresh(teamPreviewProvider(widget.code).future),
          builder: _body,
        ),
      ),
    );
  }
}
