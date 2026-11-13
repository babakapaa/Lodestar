import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/navigation/app_routes.dart';
import '../../design/components/choice_group.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import '../team/team_action.dart';
import 'class_options.dart';
import 'run_estimate.dart';
import 'run_input.dart';
import 'run_kind.dart';
import 'run_work.dart';
import 'sandbox_labels_context.dart';
import 'sandbox_providers.dart';
import 'sandbox_refresh.dart';
import 'sandbox_template.dart';

class RunLaunchScreen extends ConsumerStatefulWidget {
  const RunLaunchScreen({this.work, super.key});

  final RunWork? work;

  @override
  ConsumerState<RunLaunchScreen> createState() => _RunLaunchScreenState();
}

class _RunLaunchScreenState extends ConsumerState<RunLaunchScreen> {
  RunKind _kind = RunKind.jupyter;
  late final _title = TextEditingController(text: widget.work?.title ?? '');
  final _command = TextEditingController();
  final _hours = TextEditingController();
  String? _templateId;
  String? _teamId;
  String? _classId;
  RunEstimate? _estimate;
  Object? _estimateError;
  bool _busy = false;
  int _asked = 0;

  @override
  void dispose() {
    _title.dispose();
    _command.dispose();
    _hours.dispose();
    super.dispose();
  }

  double? get _maxHours =>
      double.tryParse(_hours.text.trim().replaceAll(',', '.'));

  Future<void> _estimateNow() async {
    final templateId = _templateId;
    if (templateId == null) {
      return;
    }
    final ask = ++_asked;
    try {
      final estimate = await ref
          .read(sandboxApiProvider)
          .estimate(
            kind: _kind,
            templateId: templateId,
            maxHours: _maxHours,
            teamId: _teamId,
          );
      if (!mounted || ask != _asked) {
        return;
      }
      setState(() {
        _estimate = estimate;
        _estimateError = null;
        final known = estimate.options.any(
          (o) => o.classId == _classId && o.enough,
        );
        if (!known) {
          _classId = estimate.suggested;
        }
      });
    } on Object catch (error) {
      if (mounted && ask == _asked) {
        setState(() => _estimateError = error);
      }
    }
  }

  void _change(VoidCallback update) {
    setState(update);
    _estimateNow();
  }

  List<String> _problems() {
    final l10n = context.l10n;
    final option = _estimate?.options
        .where((o) => o.classId == _classId)
        .firstOrNull;
    return [
      if (_title.text.trim().isEmpty) l10n.sandboxNeedTitle,
      if (_templateId == null) l10n.sandboxNeedTemplate,
      if (_kind.commanded && _command.text.trim().isEmpty)
        l10n.sandboxNeedCommand,
      if (_kind == RunKind.job && _maxHours == null) l10n.sandboxNeedHours,
      if (_templateId != null && _estimate != null && option == null)
        l10n.sandboxNeedClass,
      if (option != null && !option.enough) l10n.sandboxNotEnough,
    ];
  }

  Future<void> _start() async {
    final templateId = _templateId;
    final classId = _classId;
    if (_problems().isNotEmpty || templateId == null || classId == null) {
      setState(() {});
      return;
    }
    final work = widget.work;
    setState(() => _busy = true);
    var started = '';
    await TeamAction.run(context, () async {
      final run = await ref
          .read(sandboxApiProvider)
          .start(
            RunInput(
              kind: _kind,
              title: _title.text.trim(),
              templateId: templateId,
              classId: classId,
              command: _kind.commanded ? _command.text.trim() : '',
              maxHours: _maxHours,
              workKind: work?.kind,
              workId: work?.id,
              teamId: _teamId,
            ),
          );
      started = run.id;
    });
    if (!mounted) {
      return;
    }
    setState(() => _busy = false);
    if (started.isNotEmpty) {
      SandboxRefresh.all(ref);
      context.pushReplacement(AppRoutes.sandboxRun(started));
    }
  }

  Widget _templates(List<SandboxTemplate> all) {
    final l10n = context.l10n;
    final list = [
      for (final t in all)
        if (!t.archived) t,
    ];
    return DropdownButtonFormField<String>(
      key: const ValueKey('run-template'),
      initialValue: _templateId,
      isExpanded: true,
      decoration: InputDecoration(labelText: l10n.sandboxTemplate),
      items: [
        for (final t in list)
          DropdownMenuItem(
            value: t.id,
            child: Text(
              t.starter ? l10n.sandboxStarterName(t.name) : t.name,
              overflow: TextOverflow.ellipsis,
            ),
          ),
      ],
      onChanged: (id) => _change(() => _templateId = id),
    );
  }

  Widget _owner() {
    final l10n = context.l10n;
    final teams = ref.watch(sandboxQuotaProvider).value?.teams ?? const [];
    if (teams.isEmpty) {
      return const SizedBox.shrink();
    }
    return DropdownButtonFormField<String?>(
      key: const ValueKey('run-owner'),
      initialValue: _teamId,
      isExpanded: true,
      decoration: InputDecoration(labelText: l10n.sandboxPayFrom),
      items: [
        DropdownMenuItem(value: null, child: Text(l10n.sandboxPayMine)),
        for (final t in teams)
          DropdownMenuItem(value: t.id, child: Text(t.name)),
      ],
      onChanged: (id) => _change(() => _teamId = id),
    );
  }

  Widget _estimateView() {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    final estimate = _estimate;
    final error = _estimateError;
    if (error != null) {
      return FormError(message: ApiErrorText.of(l10n, error));
    }
    if (estimate == null) {
      return Text(l10n.sandboxPickTemplateFirst);
    }
    final muted = Theme.of(context).textTheme.bodyMedium!
        .copyWith(color: context.palette.muted);
    final night = estimate.night;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.sandboxAvailable(labels.hours(estimate.available)),
          style: muted,
        ),
        if (estimate.starter)
          Text(l10n.sandboxStarterHint, style: muted)
        else if (night != null)
          Text(
            l10n.sandboxNightHint(night.from, night.to, night.discount),
            style: muted,
          ),
        Text(l10n.sandboxRefundHint, style: muted),
        const SizedBox(height: Spacing.s),
        ClassOptions(
          estimate: estimate,
          selected: _classId,
          onSelected: (id) => setState(() => _classId = id),
        ),
        if (estimate.suggested != null)
          Align(
            alignment: Alignment.centerLeft,
            child: TextButton.icon(
              key: const ValueKey('run-suggest'),
              onPressed: () => setState(() => _classId = estimate.suggested),
              icon: const Icon(Icons.auto_awesome_outlined),
              label: Text(l10n.sandboxSuggest),
            ),
          ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final labels = context.sandboxLabels;
    final templates = ref.watch(sandboxTemplatesProvider).value ?? const [];
    final problems = _problems();
    final work = widget.work;
    return Scaffold(
      appBar: AppBar(title: Text(l10n.sandboxNewTitle)),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 720),
            child: ListView(
              padding: const EdgeInsets.all(Spacing.l),
              children: [
                if (work != null) ...[
                  Text(l10n.sandboxForWork(work.title)),
                  const SizedBox(height: Spacing.l),
                ],
                ChoiceGroup<RunKind>(
                  values: RunKind.values,
                  selected: _kind,
                  labelOf: labels.kind,
                  onSelected: (kind) => _change(() => _kind = kind),
                ),
                const SizedBox(height: Spacing.l),
                TextField(
                  key: const ValueKey('run-title'),
                  controller: _title,
                  decoration: InputDecoration(labelText: l10n.sandboxRunTitle),
                  onChanged: (_) => setState(() {}),
                ),
                _templates(templates),
                if (_kind.commanded)
                  TextField(
                    key: const ValueKey('run-command'),
                    controller: _command,
                    decoration: InputDecoration(
                      labelText: l10n.sandboxCommand,
                      helperText: _kind == RunKind.demo
                          ? l10n.sandboxDemoHint
                          : l10n.sandboxCommandHint,
                    ),
                    onChanged: (_) => setState(() {}),
                  ),
                TextField(
                  key: const ValueKey('run-hours'),
                  controller: _hours,
                  keyboardType: const TextInputType.numberWithOptions(
                    decimal: true,
                  ),
                  decoration: InputDecoration(
                    labelText: l10n.sandboxMaxHours,
                    helperText: _kind == RunKind.job
                        ? l10n.sandboxMaxHoursJob
                        : l10n.sandboxMaxHoursDefault,
                  ),
                  onChanged: (_) => _change(() {}),
                ),
                _owner(),
                const SizedBox(height: Spacing.xl),
                Section(title: l10n.sandboxCost, child: _estimateView()),
                const SizedBox(height: Spacing.xl),
                for (final p in problems)
                  ListTile(
                    dense: true,
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.info_outline, size: 18),
                    title: Text(p),
                  ),
                const SizedBox(height: Spacing.s),
                FilledButton.icon(
                  key: const ValueKey('run-start'),
                  onPressed: _busy || problems.isNotEmpty ? null : _start,
                  icon: const Icon(Icons.play_arrow),
                  label: Text(l10n.sandboxStart),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
