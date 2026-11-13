import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/files/files_providers.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'student_import.dart';
import 'student_import_report.dart';
import 'study_providers.dart';
import 'table_format.dart';

class GradebookSection extends ConsumerStatefulWidget {
  const GradebookSection({
    required this.courseId,
    required this.onImported,
    super.key,
  });

  final String courseId;
  final VoidCallback onImported;

  @override
  ConsumerState<GradebookSection> createState() => _GradebookSectionState();
}

class _GradebookSectionState extends ConsumerState<GradebookSection> {
  bool _busy = false;
  StudentImport? _report;

  Future<void> _run(Future<void> Function() action) async {
    setState(() => _busy = true);
    await TeamAction.run(context, action);
    if (mounted) {
      setState(() => _busy = false);
    }
  }

  Future<void> _export(TableFormat format) => _run(() async {
    final bytes = await ref
        .read(gradebookApiProvider)
        .export(widget.courseId, format);
    final day = DateTime.now().toIso8601String().substring(0, 10);
    await ref
        .read(exportSaverProvider)
        .save('gradebook-$day', bytes, extension: format.name);
  });

  Future<void> _import() async {
    final files = await ref.read(fileOpenerProvider).open(const [
      'csv',
      'xlsx',
    ]);
    if (files.isEmpty) {
      return;
    }
    final file = files.first;
    await _run(() async {
      final report = await ref
          .read(gradebookApiProvider)
          .importStudents(
            widget.courseId,
            file.bytes,
            TableFormat.ofName(file.name),
          );
      if (mounted) {
        setState(() => _report = report);
      }
      if (report.errors.isEmpty) {
        widget.onImported();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final report = _report;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Section(
          title: l10n.gradebookTitle,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(l10n.gradebookHint, style: muted),
              const SizedBox(height: Spacing.s),
              Wrap(
                spacing: Spacing.s,
                runSpacing: Spacing.s,
                children: [
                  for (final format in TableFormat.values)
                    OutlinedButton.icon(
                      key: ValueKey('gradebook-${format.name}'),
                      onPressed: _busy ? null : () => _export(format),
                      icon: const Icon(Icons.download_outlined),
                      label: Text(format.name.toUpperCase()),
                    ),
                ],
              ),
            ],
          ),
        ),
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.importTitle,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(l10n.importHint, style: muted),
              const SizedBox(height: Spacing.s),
              OutlinedButton.icon(
                key: const ValueKey('import-students'),
                onPressed: _busy ? null : _import,
                icon: const Icon(Icons.upload_file_outlined),
                label: Text(l10n.importAction),
              ),
              if (report != null) ...[
                const SizedBox(height: Spacing.m),
                StudentImportReport(report: report),
              ],
            ],
          ),
        ),
      ],
    );
  }
}
