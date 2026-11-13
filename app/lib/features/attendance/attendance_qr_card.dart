import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';
import 'package:qr_flutter/qr_flutter.dart';

import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'attend_link.dart';
import 'attendance_providers.dart';

class AttendanceQrCard extends ConsumerStatefulWidget {
  const AttendanceQrCard({required this.sessionId, super.key});

  final String sessionId;

  @override
  ConsumerState<AttendanceQrCard> createState() => _AttendanceQrCardState();
}

class _AttendanceQrCardState extends ConsumerState<AttendanceQrCard> {
  Timer? _timer;
  String? _link;

  @override
  void initState() {
    super.initState();
    unawaited(_load());
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  Future<void> _load() async {
    _timer?.cancel();
    try {
      final qr = await ref.read(attendanceApiProvider).qr(widget.sessionId);
      if (!mounted) {
        return;
      }
      setState(() => _link = AttendLink.build(widget.sessionId, qr.token));
      final wait = qr.expiresAt.difference(DateTime.now());
      _timer = Timer(
        wait.isNegative ? const Duration(seconds: 1) : wait,
        _load,
      );
    } on Exception {
      _timer = Timer(const Duration(seconds: 5), _load);
    }
  }

  @override
  Widget build(BuildContext context) {
    final link = _link;
    return Column(
      children: [
        if (link == null)
          const SizedBox(
            height: 280,
            child: Center(child: CircularProgressIndicator()),
          )
        else
          ColoredBox(
            color: Colors.white,
            child: QrImageView(
              key: const ValueKey('attendance-qr'),
              data: link,
              size: 280,
            ),
          ),
        const SizedBox(height: Spacing.s),
        Text(context.l10n.attendanceQrHint, textAlign: TextAlign.center),
      ],
    );
  }
}
