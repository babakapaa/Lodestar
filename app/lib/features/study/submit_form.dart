import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/offline/client_id.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'assignment_details.dart';
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
        const SizedBox(height: Spacing.l),
        Align(
          alignment: Alignment.centerLeft,
          child: LodestarButton(
            label: again ? l10n.submitAgain : l10n.submitWork,
            icon: Icons.send_outlined,
            onPressed: _busy ? null : _submit,
          ),
        ),
      ],
    );
  }
}
