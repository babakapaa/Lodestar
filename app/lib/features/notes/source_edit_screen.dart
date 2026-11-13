import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../core/offline/client_id.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../../core/sources/gost_reference.dart';
import 'notes_providers.dart';
import '../../core/sources/source.dart';
import '../../core/sources/source_kind.dart';
import 'source_labels.dart';

class SourceEditScreen extends ConsumerStatefulWidget {
  const SourceEditScreen({this.source, this.existing = false, super.key});

  final Source? source;
  final bool existing;

  @override
  ConsumerState<SourceEditScreen> createState() => _SourceEditScreenState();
}

class _SourceEditScreenState extends ConsumerState<SourceEditScreen> {
  late final String _id = widget.source?.id ?? ClientId.generate();
  late SourceKind _kind = widget.source?.kind ?? SourceKind.book;
  late final Map<String, TextEditingController> _fields = {
    for (final MapEntry(:key, :value) in _values(widget.source).entries)
      key: TextEditingController(text: value),
  };
  bool _busy = false;
  String? _error;

  static Map<String, String> _values(Source? s) => {
    'title': s?.title ?? '',
    'authors': s?.authors.join('\n') ?? '',
    'year': s?.year ?? '',
    'container': s?.container ?? '',
    'publisher': s?.publisher ?? '',
    'city': s?.city ?? '',
    'volume': s?.volume ?? '',
    'issue': s?.issue ?? '',
    'pages': s?.pages ?? '',
    'url': s?.url ?? '',
    'doi': s?.doi ?? '',
    'isbn': s?.isbn ?? '',
    'accessedAt': s?.accessedAt ?? '',
    'note': s?.note ?? '',
  };

  @override
  void dispose() {
    for (final controller in _fields.values) {
      controller.dispose();
    }
    super.dispose();
  }

  String _text(String key) => _fields[key]!.text.trim();

  Source _source() => Source(
    id: _id,
    kind: _kind,
    title: _text('title'),
    authors: [
      for (final line in _fields['authors']!.text.split('\n'))
        if (line.trim().isNotEmpty) line.trim(),
    ],
    year: _text('year'),
    container: _text('container'),
    publisher: _text('publisher'),
    city: _text('city'),
    volume: _text('volume'),
    issue: _text('issue'),
    pages: _text('pages'),
    url: _text('url'),
    doi: _text('doi'),
    isbn: _text('isbn'),
    accessedAt: _text('accessedAt'),
    note: _fields['note']!.text,
  );

  Future<void> _run(Future<void> Function() action) async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await action();
      ref.invalidate(sourcesProvider);
      if (mounted) {
        context.pop();
      }
    } on ApiException catch (error) {
      setState(() {
        _busy = false;
        _error = ApiErrorText.of(l10n, error);
      });
    }
  }

  Future<void> _copy() async {
    final messenger = ScaffoldMessenger.of(context);
    final copied = context.l10n.sourcesCopiedOne;
    await Clipboard.setData(
      ClipboardData(text: GostReference.format(_source())),
    );
    messenger.showSnackBar(SnackBar(content: Text(copied)));
  }

  Widget _field(String key, String label, {int lines = 1}) => Padding(
    key: ValueKey(key),
    padding: const EdgeInsets.only(bottom: Spacing.m),
    child: TextField(
      controller: _fields[key],
      enabled: !_busy,
      minLines: lines,
      maxLines: lines == 1 ? 1 : null,
      decoration: InputDecoration(labelText: label),
      onChanged: (_) => setState(() {}),
    ),
  );

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = SourceLabels(l10n);
    final error = _error;
    final preview = _text('title').isEmpty
        ? null
        : GostReference.format(_source());
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.existing ? l10n.sourcesEdit : l10n.sourcesNew),
        actions: [
          if (preview != null)
            IconButton(
              tooltip: l10n.sourcesCopyOne,
              icon: const Icon(Icons.copy),
              onPressed: _copy,
            ),
          if (widget.existing)
            IconButton(
              tooltip: l10n.sourcesDelete,
              icon: const Icon(Icons.delete_outline),
              onPressed: _busy
                  ? null
                  : () => _run(() => ref.read(sourcesApiProvider).delete(_id)),
            ),
        ],
      ),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 720),
          child: ListView(
            padding: const EdgeInsets.all(Spacing.l),
            children: [
              if (preview != null) ...[
                SelectableText(preview),
                const SizedBox(height: Spacing.l),
              ],
              Padding(
                key: const ValueKey('kind'),
                padding: const EdgeInsets.only(bottom: Spacing.m),
                child: DropdownButtonFormField<SourceKind>(
                  initialValue: _kind,
                  isExpanded: true,
                  decoration: InputDecoration(labelText: l10n.sourcesKind),
                  items: [
                    for (final kind in SourceKind.values)
                      DropdownMenuItem(
                        value: kind,
                        child: Text(labels.kind(kind)),
                      ),
                  ],
                  onChanged: _busy
                      ? null
                      : (kind) => setState(() => _kind = kind ?? _kind),
                ),
              ),
              _field('title', l10n.sourcesFieldTitle),
              _field('authors', l10n.sourcesFieldAuthors, lines: 2),
              _field('year', l10n.sourcesFieldYear),
              _field('container', l10n.sourcesFieldContainer),
              _field('publisher', l10n.sourcesFieldPublisher),
              _field('city', l10n.sourcesFieldCity),
              _field('volume', l10n.sourcesFieldVolume),
              _field('issue', l10n.sourcesFieldIssue),
              _field('pages', l10n.sourcesFieldPages),
              _field('url', l10n.sourcesFieldUrl),
              _field('accessedAt', l10n.sourcesFieldAccessed),
              _field('doi', 'DOI'),
              _field('isbn', 'ISBN'),
              _field('note', l10n.sourcesFieldNote, lines: 3),
              if (error != null) ...[
                FormError(message: error),
                const SizedBox(height: Spacing.m),
              ],
              Align(
                alignment: Alignment.centerLeft,
                child: LodestarButton(
                  label: l10n.actionSave,
                  icon: Icons.check,
                  onPressed: _busy || _text('title').isEmpty
                      ? null
                      : () => _run(
                          () => ref.read(sourcesApiProvider).save(_source()),
                        ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
