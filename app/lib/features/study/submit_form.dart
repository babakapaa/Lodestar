import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/files/files_providers.dart';
import '../../core/offline/client_id.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_details.dart';
import 'attached_file.dart';
import 'search_map.dart';
import 'study_providers.dart';
import 'submission_draft.dart';

class SubmitForm extends ConsumerStatefulWidget {
  const SubmitForm({required this.details, super.key});

  final AssignmentDetails details;

  @override
  ConsumerState<SubmitForm> createState() => _SubmitFormState();
}

class _SubmitFormState extends ConsumerState<SubmitForm> {
  static const maxFiles = 20;

  final _body = TextEditingController();
  final _files = <AttachedFile>[];
  String? _error;
  bool _busy = false;
  bool _attaching = false;

  @override
  void dispose() {
    _body.dispose();
    super.dispose();
  }

  SearchMap? get _map {
    final map = ref
        .read(searchMapProvider(widget.details.assignment.id))
        .value
        ?.value
        .map;
    return map == null ||
            (map.isEmpty && widget.details.assignment.searchMap == null)
        ? null
        : map;
  }

  Future<void> _attach() async {
    final l10n = context.l10n;
    setState(() {
      _attaching = true;
      _error = null;
    });
    try {
      final opened = await ref.read(fileOpenerProvider).open(const []);
      final transfer = ref.read(fileTransferProvider);
      for (final file in opened.take(maxFiles - _files.length)) {
        final uploaded = await transfer.upload(file);
        if (!mounted) {
          return;
        }
        setState(
          () => _files.add(
            AttachedFile(
              id: uploaded.id,
              name: uploaded.name,
              size: uploaded.size,
            ),
          ),
        );
      }
    } on Object catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    } finally {
      if (mounted) {
        setState(() => _attaching = false);
      }
    }
  }

  Future<void> _submit() async {
    final l10n = context.l10n;
    final body = _body.text.trim();
    if (body.isEmpty && _files.isEmpty) {
      setState(() => _error = l10n.errorInvalidSubmission);
      return;
    }
    final assignment = widget.details.assignment;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      await ref
          .read(outboxProvider.notifier)
          .add(
            SubmissionDraft(
              clientId: ClientId.generate(),
              assignmentId: assignment.id,
              assignmentTitle: assignment.title,
              body: body,
              composedAt: DateTime.now(),
              files: List.of(_files),
              searchMap: _map,
            ),
          );
      _body.clear();
      _files.clear();
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final error = _error;
    final again = widget.details.mine?.latest != null;
    final rules = widget.details.assignment.searchMap;
    final map = ref
        .watch(searchMapProvider(widget.details.assignment.id))
        .value
        ?.value
        .map;
    final blocked =
        rules != null && (map == null || map.gaps(rules).isNotEmpty);
    final attached = !blocked && map != null && (rules != null || !map.isEmpty);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TextField(
          controller: _body,
          minLines: 4,
          maxLines: 12,
          decoration: InputDecoration(
            labelText: l10n.submitHint,
            alignLabelWithHint: true,
          ),
        ),
        if (_files.isNotEmpty) ...[
          const SizedBox(height: Spacing.s),
          Wrap(
            spacing: Spacing.s,
            runSpacing: Spacing.xs,
            children: [
              for (final file in _files)
                InputChip(
                  avatar: const Icon(Icons.attach_file, size: 16),
                  label: Text(file.name, overflow: TextOverflow.ellipsis),
                  onDeleted: () => setState(() => _files.remove(file)),
                ),
            ],
          ),
        ],
        if (error != null) FormError(message: error),
        if (blocked) ...[
          const SizedBox(height: Spacing.m),
          NoticeBanner(icon: Icons.travel_explore, message: l10n.mapNotReady),
        ] else if (attached) ...[
          const SizedBox(height: Spacing.m),
          Text(l10n.mapAttached),
        ],
        const SizedBox(height: Spacing.l),
        Wrap(
          spacing: Spacing.m,
          runSpacing: Spacing.s,
          children: [
            LodestarButton(
              label: again ? l10n.submitAgain : l10n.submitWork,
              icon: Icons.send_outlined,
              onPressed: _busy || blocked || _attaching ? null : _submit,
            ),
            TextButton.icon(
              key: const ValueKey('submit-attach'),
              onPressed: _attaching || _files.length >= maxFiles
                  ? null
                  : _attach,
              icon: _attaching
                  ? const SizedBox.square(
                      dimension: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.attach_file),
              label: Text(l10n.submitAttach),
            ),
          ],
        ),
      ],
    );
  }
}
