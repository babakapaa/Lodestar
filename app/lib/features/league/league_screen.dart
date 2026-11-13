import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'league_lead_body.dart';
import 'league_providers.dart';
import 'league_student_body.dart';
import 'league_view.dart';

class LeagueScreen extends ConsumerWidget {
  const LeagueScreen({required this.kind, required this.originId, super.key});

  final String kind;
  final String originId;

  List<Widget> _body(BuildContext context, LeagueView view) {
    final textTheme = Theme.of(context).textTheme;
    return [
      Text(view.title, style: textTheme.headlineSmall),
      Text(
        context.l10n.leagueSeason(view.season.number),
        style: textTheme.bodyMedium!.copyWith(color: context.palette.muted),
      ),
      const SizedBox(height: Spacing.xl),
      if (view.leads)
        LeagueLeadBody(view: view)
      else
        LeagueStudentBody(view: view),
    ];
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final key = (kind, originId);
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.leagueTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(leagueProvider(key)),
          onRefresh: () => ref.refresh(leagueProvider(key).future),
          builder: (view) => _body(context, view),
        ),
      ),
    );
  }
}
