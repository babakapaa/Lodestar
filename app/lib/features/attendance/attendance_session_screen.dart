import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:material_ui/material_ui.dart';

import '../../design/components/async_list_view.dart';
import '../../design/components/lodestar_button.dart';
import '../../design/components/lodestar_button_variant.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/study_labels_context.dart';
import '../team/team_action.dart';
import 'attendance_providers.dart';
import 'attendance_qr_card.dart';
import 'attendance_session.dart';
import 'attendance_status.dart';
import 'mark_menu.dart';

class AttendanceSessionScreen extends ConsumerStatefulWidget {
  const AttendanceSessionScreen({required this.sessionId, super.key});

  final String sessionId;

  @override
  ConsumerState<AttendanceSessionScreen> createState() =>
      _AttendanceSessionScreenState();
}

class _AttendanceSessionScreenState
    extends ConsumerState<AttendanceSessionScreen> {
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _timer = Timer.periodic(const Duration(seconds: 10), (_) {
      if (ref.read(attendanceSessionProvider(widget.sessionId)).value?.open ??
          false) {
        ref.invalidate(attendanceSessionProvider(widget.sessionId));
      }
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  Future<void> _act(Future<Object?> Function() action) async {
    await TeamAction.run(context, () async {
      await action();
    });
    ref.invalidate(attendanceSessionProvider(widget.sessionId));
  }

  void _mark(String userId, AttendanceStatus? status) => _act(
    () =>
        ref.read(attendanceApiProvider).mark(widget.sessionId, userId, status),
  );

  List<Widget> _body(AttendanceSession session) {
    final l10n = context.l10n;
    final labels = context.studyLabels;
    final textTheme = Theme.of(context).textTheme;
    final api = ref.read(attendanceApiProvider);
    final until = session.codeUntil;
    return [
      Text(session.courseTitle, style: textTheme.titleLarge),
      Text(
        [
          labels.day(session.startsAt),
          '${labels.clock(session.startsAt)}–${labels.clock(session.endsAt)}',
          if (session.room.isNotEmpty) session.room,
        ].join(', '),
      ),
      const SizedBox(height: Spacing.l),
      if (!session.open)
        NoticeBanner(icon: Icons.lock_outline, message: l10n.attendanceClosed),
      if (session.open && session.methods.qr) ...[
        LodestarCard(child: AttendanceQrCard(sessionId: session.id)),
        const SizedBox(height: Spacing.l),
      ],
      if (session.open && session.code.isNotEmpty) ...[
        LodestarCard(
          padding: EdgeInsets.zero,
          child: ListTile(
            contentPadding: LodestarCard.tilePadding,
            title: Text(
              l10n.attendanceCode(session.code),
              key: const ValueKey('attendance-code'),
              style: textTheme.headlineMedium,
            ),
            subtitle: until == null
                ? null
                : Text(l10n.attendanceCodeUntil(labels.clock(until))),
            trailing: TextButton.icon(
              key: const ValueKey('attendance-new-code'),
              onPressed: () => _act(() => api.refreshCode(session.id)),
              icon: const Icon(Icons.refresh),
              label: Text(l10n.attendanceNewCode),
            ),
          ),
        ),
        const SizedBox(height: Spacing.l),
      ],
      Section(
        title: l10n.attendanceHere(session.counts.here, session.marks.length),
        child: Column(
          children: [
            for (final mark in session.marks)
              ListTile(
                key: ValueKey('student-${mark.userId}'),
                contentPadding: EdgeInsets.zero,
                title: Text(mark.name),
                trailing: MarkMenu(
                  mark: mark,
                  onChanged: (status) => _mark(mark.userId, status),
                ),
              ),
          ],
        ),
      ),
      if (session.open) ...[
        const SizedBox(height: Spacing.xl),
        Text(l10n.attendanceCloseHint),
        const SizedBox(height: Spacing.s),
        LodestarButton(
          key: const ValueKey('attendance-close'),
          label: l10n.attendanceClose,
          icon: Icons.lock_outline,
          variant: LodestarButtonVariant.secondary,
          onPressed: () => _act(() => api.close(session.id)),
        ),
      ],
    ];
  }

  @override
  Widget build(BuildContext context) {
    final session = ref.watch(attendanceSessionProvider(widget.sessionId));
    return Scaffold(
      appBar: AppBar(title: Text(context.l10n.attendanceSessionTitle)),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: session,
          onRefresh: () =>
              ref.refresh(attendanceSessionProvider(widget.sessionId).future),
          builder: _body,
        ),
      ),
    );
  }
}
