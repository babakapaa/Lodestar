import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/files/files_providers.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'map_requirements.dart';
import 'map_rules.dart';
import 'map_source.dart';
import 'map_source_card.dart';
import 'map_source_picker.dart';
import 'search_map.dart';
import 'search_map_diagram.dart';
import 'search_map_draft.dart';
import 'search_map_markdown.dart';
import 'study_providers.dart';

class SearchMapEditor extends ConsumerStatefulWidget {
  const SearchMapEditor({
    required this.assignmentId,
    required this.title,
    required this.draft,
    super.key,
  });

  final String assignmentId;
  final String title;
  final SearchMapDraft draft;

  @override
  ConsumerState<SearchMapEditor> createState() => _SearchMapEditorState();
}

class _SearchMapEditorState extends ConsumerState<SearchMapEditor> {
  late SearchMap _map = widget.draft.map;
  late final _path = TextEditingController(text: _map.path);
  late final _queries = TextEditingController(text: _map.queries.join('\n'));
  bool _busy = false;
  String? _error;

  MapRules? get _rules => widget.draft.rules;

  @override
  void dispose() {
    _path.dispose();
    _queries.dispose();
    super.dispose();
  }

  SearchMap get _current => _map.copyWith(
    path: _path.text,
    queries: [
      for (final line in _queries.text.split('\n'))
        if (line.trim().isNotEmpty) line.trim(),
    ],
  );

  Future<void> _add() async {
    final picked = await MapSourcePicker.show(context, _map.contains);
    if (picked == null) {
      return;
    }
    final (source, used) = picked;
    final entry = MapSource.fromSource(source);
    setState(
      () => _map = used
          ? _map.copyWith(used: [..._map.used, entry])
          : _map.copyWith(rejected: [..._map.rejected, entry]),
    );
  }

  void _update(bool rejected, int index, MapSource source) {
    final list = [...(rejected ? _map.rejected : _map.used)]..[index] = source;
    setState(
      () => _map = rejected
          ? _map.copyWith(rejected: list)
          : _map.copyWith(used: list),
    );
  }

  void _move(bool rejected, int index) {
    final from = [...(rejected ? _map.rejected : _map.used)];
    final source = from.removeAt(index);
    final to = [...(rejected ? _map.used : _map.rejected), source];
    setState(
      () => _map = rejected
          ? _map.copyWith(rejected: from, used: to)
          : _map.copyWith(used: from, rejected: to),
    );
  }

  void _remove(bool rejected, int index) {
    final list = [...(rejected ? _map.rejected : _map.used)]..removeAt(index);
    setState(
      () => _map = rejected
          ? _map.copyWith(rejected: list)
          : _map.copyWith(used: list),
    );
  }

  Future<void> _save() async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref
          .read(studyRepositoryProvider)
          .saveSearchMap(widget.assignmentId, _current);
      ref.invalidate(searchMapProvider(widget.assignmentId));
      messenger.showSnackBar(SnackBar(content: Text(l10n.mapSaved)));
    } on ApiException catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  Future<void> _export() async {
    final text = SearchMapMarkdown.of(context.l10n, widget.title, _current);
    await ref
        .read(exportSaverProvider)
        .save('search-map', utf8.encode(text), extension: 'md');
  }

  List<Widget> _sources(bool rejected) {
    final l10n = context.l10n;
    final list = rejected ? _map.rejected : _map.used;
    return [
      if (list.isEmpty)
        Text(l10n.mapNoneHere, style: TextStyle(color: context.palette.muted)),
      for (final (index, source) in list.indexed)
        Padding(
          key: ValueKey(('$rejected', source.sourceId, source.reference)),
          padding: const EdgeInsets.only(bottom: Spacing.m),
          child: MapSourceCard(
            source: source,
            rejected: rejected,
            onChanged: (next) => _update(rejected, index, next),
            onMove: () => _move(rejected, index),
            onRemove: () => _remove(rejected, index),
          ),
        ),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    final current = _current;
    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.mapTitle),
        actions: [
          IconButton(
            tooltip: l10n.mapExport,
            icon: const Icon(Icons.download_outlined),
            onPressed: _export,
          ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 720),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                if (widget.title.isNotEmpty)
                  Text(
                    widget.title,
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                const SizedBox(height: Spacing.s),
                MapRequirements(map: current, rules: _rules),
                const SizedBox(height: Spacing.s),
                Text(
                  l10n.mapFacts,
                  style: TextStyle(color: context.palette.muted),
                ),
                const SizedBox(height: Spacing.xl),
                Section(
                  title: l10n.mapUsed,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: _sources(false),
                  ),
                ),
                const SizedBox(height: Spacing.l),
                Section(
                  title: l10n.mapRejected,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: _sources(true),
                  ),
                ),
                const SizedBox(height: Spacing.m),
                Align(
                  alignment: Alignment.centerLeft,
                  child: LodestarButton(
                    label: l10n.mapAddSource,
                    icon: Icons.add,
                    variant: LodestarButtonVariant.secondary,
                    onPressed: _add,
                  ),
                ),
                const SizedBox(height: Spacing.xl),
                Section(
                  title: l10n.mapPath,
                  child: Column(
                    children: [
                      TextField(
                        controller: _path,
                        minLines: 3,
                        maxLines: 10,
                        decoration: InputDecoration(
                          labelText: l10n.mapPathHint,
                          alignLabelWithHint: true,
                        ),
                        onChanged: (_) => setState(() {}),
                      ),
                      const SizedBox(height: Spacing.m),
                      TextField(
                        controller: _queries,
                        minLines: 2,
                        maxLines: 8,
                        decoration: InputDecoration(
                          labelText: l10n.mapQueries,
                          alignLabelWithHint: true,
                        ),
                        onChanged: (_) => setState(() {}),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: Spacing.xl),
                Section(
                  title: l10n.mapScheme,
                  child: SearchMapDiagram(map: current),
                ),
                if (error != null) FormError(message: error),
                const SizedBox(height: Spacing.l),
                Align(
                  alignment: Alignment.centerLeft,
                  child: LodestarButton(
                    label: l10n.mapSave,
                    icon: Icons.check,
                    onPressed: _busy ? null : _save,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
