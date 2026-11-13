import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import 'league_providers.dart';
import 'league_rules_form.dart';

class LeagueRulesScreen extends ConsumerWidget {
  const LeagueRulesScreen({
    required this.kind,
    required this.originId,
    super.key,
  });

  final String kind;
  final String originId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final key = (kind, originId);
    return Scaffold(
      appBar: AppBar(),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: ref.watch(leagueProvider(key)),
          onRefresh: () => ref.refresh(leagueProvider(key).future),
          builder: (view) => [LeagueRulesForm(view: view)],
        ),
      ),
    );
  }
}
