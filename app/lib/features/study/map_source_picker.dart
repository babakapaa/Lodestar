import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/sources/gost_reference.dart';
import '../../core/sources/source.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'study_providers.dart';

class MapSourcePicker extends ConsumerWidget {
  const MapSourcePicker({required this.taken, super.key});

  final bool Function(String sourceId) taken;

  static Future<(Source, bool)?> show(
    BuildContext context,
    bool Function(String sourceId) taken,
  ) => showModalBottomSheet<(Source, bool)>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (context) => MapSourcePicker(taken: taken),
  );

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = context.l10n;
    final loaded = ref.watch(mapSourcesProvider);
    final list = [
      for (final source in loaded.value?.value ?? const <Source>[])
        if (!taken(source.id)) source,
    ];
    return SafeArea(
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxHeight: MediaQuery.sizeOf(context).height * 0.75,
        ),
        child: loaded.value == null
            ? const Padding(
                padding: EdgeInsets.all(Spacing.xl),
                child: Center(child: CircularProgressIndicator()),
              )
            : list.isEmpty
            ? Padding(
                padding: const EdgeInsets.all(Spacing.l),
                child: Text(l10n.mapNoSources),
              )
            : ListView(
                shrinkWrap: true,
                padding: const EdgeInsets.symmetric(horizontal: Spacing.l),
                children: [
                  for (final source in list)
                    Padding(
                      padding: const EdgeInsets.only(bottom: Spacing.m),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(GostReference.format(source)),
                          Wrap(
                            spacing: Spacing.s,
                            children: [
                              TextButton(
                                onPressed: () =>
                                    Navigator.of(context).pop((source, true)),
                                child: Text(l10n.mapPickUsed),
                              ),
                              TextButton(
                                onPressed: () =>
                                    Navigator.of(context).pop((source, false)),
                                child: Text(l10n.mapPickRejected),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                ],
              ),
      ),
    );
  }
}
