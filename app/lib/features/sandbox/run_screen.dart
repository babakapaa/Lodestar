import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/navigation/app_routes.dart';
import '../../design/components/async_list_view.dart';
import '../../design/components/lodestar_card.dart';
import '../../design/components/notice_banner.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../profile/profile_providers.dart';
import '../study/study_labels_context.dart';
import '../team/team_action.dart';
import 'hours_dialog.dart';
import 'run_link.dart';
import 'run_status.dart';
import 'sandbox_labels_context.dart';
import 'sandbox_providers.dart';
import 'sandbox_refresh.dart';
import 'sandbox_run.dart';

class RunScreen extends ConsumerStatefulWidget {
  const RunScreen({required this.runId, super.key});

  final String runId;

  @override
  ConsumerState<RunScreen> createState() => _RunScreenState();
}

class _RunScreenState extends ConsumerState<RunScreen> {
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _timer = Timer.periodic(const Duration(seconds: 15), (_) {
      final run = ref.read(sandboxRunProvider(widget.runId)).value;
      if (run != null && run.status.active) {
        ref.invalidate(sandboxRunProvider(widget.runId));
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
    SandboxRefresh.all(ref);
  }

  Future<void> _rerun(SandboxRun run) async {
    var next = '';
    await TeamAction.run(context, () async {
      next = (await ref.read(sandboxApiProvider).rerun(run.id)).id;
    });
    SandboxRefresh.all(ref);
    if (next.isNotEmpty && mounted) {
      context.pushReplacement(AppRoutes.sandboxRun(next));
    }
  }

  Future<void> _askHours(SandboxRun run) async {
    final l10n = context.l10n;
    final request = await HoursDialog.show(
      context,
      title: l10n.sandboxAskTitle,
      action: l10n.sandboxAsk,
      withReason: true,
    );
    final kind = run.workKind;
    final id = run.workId;
    if (request == null || kind == null || id == null || !mounted) {
      return;
    }
    final sent = await TeamAction.run(
      context,
      () => ref
          .read(sandboxApiProvider)
          .requestTopup(
            workKind: kind,
            workId: id,
            hours: request.hours,
            reason: request.reason,
          ),
    );
    if (sent && mounted) {
      ScaffoldMessenger.of(context)
          .showSnackBar(SnackBar(content: Text(l10n.sandboxAskSent)));
    }
  }

  List<Widget> _notices(SandboxRun run) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    return [
      if (run.status == RunStatus.held)
        NoticeBanner(
          icon: Icons.pause_circle_outline,
          message: l10n.sandboxHeldNotice(run.holdNote),
        )
      else if (labels.reason(run.reason) case final reason?)
        NoticeBanner(icon: Icons.info_outline, message: reason),
      if (run.refunded)
        NoticeBanner(icon: Icons.undo, message: l10n.sandboxRefunded),
      if (run.wasteful)
        NoticeBanner(
          icon: Icons.savings_outlined,
          message: l10n.sandboxSaveTip(run.gpuLoad!.round()),
        ),
    ];
  }

  List<Widget> _body(SandboxRun run) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    final study = context.studyLabels;
    final me = ref.watch(accountProvider).value?.id;
    final mine = me != null && me == run.userId;
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final facts = <(String, String)>[
      (l10n.sandboxFactKind, labels.kind(run.kind)),
      (l10n.sandboxFactClass, run.className),
      (l10n.sandboxTemplate, run.templateName),
      if (run.teamName.isNotEmpty) (l10n.sandboxFactTeam, run.teamName),
      if (run.workTitle.isNotEmpty) (l10n.sandboxFactWork, run.workTitle),
      if (!mine) (l10n.sandboxFactBy, run.displayName),
      (l10n.sandboxMaxHours, labels.hours(run.maxHours)),
      if (run.free)
        (l10n.sandboxFactCharged, l10n.sandboxFree)
      else if (run.charged case final charged?)
        (l10n.sandboxFactCharged, labels.hours(charged)),
      if (run.startedAt case final at?)
        (l10n.sandboxFactStarted, study.time(at)),
      if (run.endedAt case final at?) (l10n.sandboxFactEnded, study.time(at)),
      if (run.exitCode case final code?) (l10n.sandboxFactExit, '$code'),
      if (run.command.isNotEmpty) (l10n.sandboxCommand, run.command),
    ];
    return [
      Text(
        labels.summary(run),
        key: const ValueKey('run-status'),
        style: Theme.of(context).textTheme.titleMedium,
      ),
      const SizedBox(height: Spacing.m),
      for (final notice in _notices(run)) ...[
        notice,
        const SizedBox(height: Spacing.s),
      ],
      if (mine && run.url.isNotEmpty && run.status == RunStatus.running) ...[
        RunLink(run: run),
        const SizedBox(height: Spacing.l),
      ],
      LodestarCard(
        child: Column(
          children: [
            for (final (name, value) in facts)
              ListTile(
                dense: true,
                contentPadding: EdgeInsets.zero,
                title: Text(name, style: muted),
                subtitle: SelectableText(value),
              ),
          ],
        ),
      ),
      if (mine) ...[
        const SizedBox(height: Spacing.l),
        Wrap(
          spacing: Spacing.s,
          runSpacing: Spacing.s,
          children: [
            if (run.status.active && run.status != RunStatus.stopping)
              OutlinedButton.icon(
                key: const ValueKey('run-stop'),
                onPressed: () =>
                    _act(() => ref.read(sandboxApiProvider).stop(run.id)),
                icon: const Icon(Icons.stop_circle_outlined),
                label: Text(l10n.sandboxStop),
              ),
            if (!run.status.active)
              FilledButton.icon(
                key: const ValueKey('run-rerun'),
                onPressed: () => _rerun(run),
                icon: const Icon(Icons.replay),
                label: Text(l10n.sandboxRerun),
              ),
            if (run.workKind != null)
              TextButton.icon(
                key: const ValueKey('run-ask'),
                onPressed: () => _askHours(run),
                icon: const Icon(Icons.volunteer_activism_outlined),
                label: Text(l10n.sandboxAsk),
              ),
          ],
        ),
      ],
      if (run.log.isNotEmpty) ...[
        const SizedBox(height: Spacing.xl),
        Section(
          title: l10n.sandboxLog,
          child: LodestarCard(
            child: SelectableText(
              run.log,
              key: const ValueKey('run-log'),
              style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
            ),
          ),
        ),
      ],
    ];
  }

  @override
  Widget build(BuildContext context) {
    final run = ref.watch(sandboxRunProvider(widget.runId));
    return Scaffold(
      appBar: AppBar(
        title: Text(run.value?.title ?? context.l10n.sandboxTitle),
      ),
      body: SafeArea(
        top: false,
        child: AsyncListView(
          value: run,
          onRefresh: () => ref.refresh(sandboxRunProvider(widget.runId).future),
          builder: _body,
        ),
      ),
    );
  }
}
