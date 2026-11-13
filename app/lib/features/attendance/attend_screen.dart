import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/form_error.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'attendance_check.dart';
import 'attendance_labels_context.dart';
import 'attendance_providers.dart';

class AttendScreen extends ConsumerStatefulWidget {
  const AttendScreen({this.sessionId, this.token, super.key});

  final String? sessionId;
  final String? token;

  @override
  ConsumerState<AttendScreen> createState() => _AttendScreenState();
}

class _AttendScreenState extends ConsumerState<AttendScreen> {
  final _code = TextEditingController();
  AttendanceCheck? _done;
  String? _error;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    final session = widget.sessionId;
    final token = widget.token;
    if (session != null && token != null) {
      WidgetsBinding.instance.addPostFrameCallback(
        (_) => unawaited(_check(sessionId: session, token: token)),
      );
    }
  }

  @override
  void dispose() {
    _code.dispose();
    super.dispose();
  }

  Future<void> _check({String? sessionId, String? token, String? code}) async {
    final l10n = context.l10n;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      final done = await ref
          .read(attendanceApiProvider)
          .check(sessionId: sessionId, token: token, code: code);
      ref.invalidate(openLessonsProvider);
      setState(() => _done = done);
    } on ApiException catch (error) {
      setState(() => _error = ApiErrorText.of(l10n, error));
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final done = _done;
    final error = _error;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.attendanceCheckTitle)),
      body: SafeArea(
        top: false,
        child: ListView(
          padding: const EdgeInsets.all(Spacing.l),
          children: [
            if (done != null)
              NoticeBanner(
                key: const ValueKey('attendance-done'),
                icon: Icons.check_circle_outline,
                message: l10n.attendanceDone(
                  done.courseTitle,
                  context.attendanceLabels.status(done.status),
                ),
              )
            else ...[
              Text(l10n.attendanceScanHint),
              const SizedBox(height: Spacing.l),
              TextField(
                key: const ValueKey('attendance-code-field'),
                controller: _code,
                keyboardType: TextInputType.number,
                maxLength: 4,
                decoration: InputDecoration(
                  labelText: l10n.attendanceCodeField,
                ),
                onSubmitted: (_) => _check(code: _code.text.trim()),
              ),
              if (error != null) FormError(message: error),
              const SizedBox(height: Spacing.l),
              LodestarButton(
                key: const ValueKey('attendance-check'),
                label: _busy ? l10n.attendanceChecking : l10n.attendanceCheckIn,
                icon: Icons.how_to_reg_outlined,
                onPressed: _busy ? null : () => _check(code: _code.text.trim()),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
