import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/offline/client_id.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_details.dart';
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
  final _body = TextEditingController();
  String? _error;
  bool _busy = false;

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

  Future<void> _submit() async {
    final l10n = context.l10n;
    final body = _body.text.trim();
    if (body.isEmpty) {
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
              searchMap: _map,
            ),
          );
      _body.clear();
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
        if (error != null) FormError(message: error),
        if (blocked) ...[
          const SizedBox(height: Spacing.m),
          NoticeBanner(icon: Icons.travel_explore, message: l10n.mapNotReady),
        ] else if (attached) ...[
          const SizedBox(height: Spacing.m),
          Text(l10n.mapAttached),
        ],
        const SizedBox(height: Spacing.l),
        Align(
          alignment: Alignment.centerLeft,
          child: LodestarButton(
            label: again ? l10n.submitAgain : l10n.submitWork,
            icon: Icons.send_outlined,
            onPressed: _busy || blocked ? null : _submit,
          ),
        ),
      ],
    );
  }
}
