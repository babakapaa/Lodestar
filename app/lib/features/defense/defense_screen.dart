import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../design/components/form_error.dart';
import '../../design/components/section.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import '../study/study_providers.dart';
import 'defense.dart';
import 'defense_labels.dart';
import 'defense_providers.dart';
import 'defense_question.dart';
import 'defense_row.dart';
import 'finish_defense_dialog.dart';
import 'new_question_dialog.dart';
import 'question_card.dart';

class DefenseScreen extends ConsumerStatefulWidget {
  const DefenseScreen({
    required this.assignmentId,
    required this.studentId,
    super.key,
  });

  final String assignmentId;
  final String studentId;

  @override
  ConsumerState<DefenseScreen> createState() => _DefenseScreenState();
}

class _DefenseScreenState extends ConsumerState<DefenseScreen> {
  final _summary = TextEditingController();
  List<DefenseQuestion>? _questions;
  Defense? _defense;
  Object? _error;
  bool _busy = false;
  int _revision = 0;

  String get _a => widget.assignmentId;

  String get _s => widget.studentId;

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _summary.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final d = await ref.read(defenseApiProvider).get(_a, _s);
      if (mounted) {
        setState(() {
          _defense = d;
          _questions = List.of(d.questions);
          _summary.text = d.summary;
        });
      }
    } on Object catch (error) {
      if (mounted) {
        setState(() => _error = error);
      }
    }
  }

  Future<bool> _run(Future<void> Function() request) async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    setState(() => _busy = true);
    try {
      await request();
      return true;
    } on Object catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
      return false;
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  Future<void> _add() async {
    final question = await NewQuestionDialog.show(context);
    if (question != null) {
      setState(() {
        _questions = [...?_questions, question];
        _revision++;
      });
    }
  }

  Future<void> _draftProtocol() => _run(() async {
    final summary = await ref
        .read(defenseApiProvider)
        .draftProtocol(_a, _s, _questions ?? const [], _summary.text.trim());
    _summary.text = summary;
  });

  Future<bool> _save() => _run(() async {
    final saved = await ref
        .read(defenseApiProvider)
        .save(_a, _s, _questions ?? const [], _summary.text.trim());
    setState(() => _defense = saved);
  });

  Future<void> _finish() async {
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    final router = GoRouter.of(context);
    final assignment = ref.read(assignmentProvider(_a)).value?.value.assignment;
    if (assignment == null) {
      return;
    }
    final choice = await FinishDefenseDialog.show(context, assignment);
    if (choice == null || !await _save()) {
      return;
    }
    final done = await _run(
      () => ref
          .read(defenseApiProvider)
          .finish(_a, _s, choice.status, choice.grade),
    );
    if (!done) {
      return;
    }
    ref
      ..invalidate(defenseRowsProvider(_a))
      ..invalidate(assignmentWorkProvider(_a));
    messenger.showSnackBar(SnackBar(content: Text(l10n.defenseFinished)));
    router.pop();
  }

  void _change(int index, DefenseQuestion question) {
    _questions![index] = question;
    setState(() {});
  }

  void _remove(int index) => setState(() {
    _questions!.removeAt(index);
    _revision++;
  });

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final rows =
        ref.watch(defenseRowsProvider(_a)).value ?? const <DefenseRow>[];
    ref.watch(assignmentProvider(_a));
    final name = [
      for (final row in rows)
        if (row.student.userId == _s) row.student.label,
    ].firstOrNull;
    final questions = _questions;
    final defense = _defense;
    return Scaffold(
      appBar: AppBar(title: Text(name ?? l10n.defenseTitle)),
      body: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 760),
          child: _error != null
              ? Padding(
                  padding: const EdgeInsets.all(Spacing.l),
                  child: FormError(message: ApiErrorText.of(l10n, _error!)),
                )
              : questions == null || defense == null
              ? const Center(child: CircularProgressIndicator())
              : ListView(
                  padding: const EdgeInsets.all(Spacing.l),
                  children: [
                    if (defense.status.finished) ...[
                      Text(DefenseLabels(l10n).status(defense.status)),
                      const SizedBox(height: Spacing.l),
                    ],
                    Section(
                      title: l10n.defenseQuestions,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          if (questions.isEmpty) Text(l10n.defenseNoQuestions),
                          for (final (i, q) in questions.indexed) ...[
                            QuestionCard(
                              key: ValueKey('q-$_revision-$i'),
                              index: i,
                              question: q,
                              onChanged: (changed) => _change(i, changed),
                              onRemove: () => _remove(i),
                            ),
                            const SizedBox(height: Spacing.m),
                          ],
                          Align(
                            alignment: Alignment.centerLeft,
                            child: OutlinedButton.icon(
                              key: const ValueKey('add-question'),
                              onPressed: _busy ? null : _add,
                              icon: const Icon(Icons.add),
                              label: Text(l10n.defenseAddQuestion),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: Spacing.xl),
                    Section(
                      title: l10n.defenseProtocol,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          TextField(
                            key: const ValueKey('defense-summary'),
                            controller: _summary,
                            minLines: 3,
                            maxLines: 10,
                            maxLength: 20000,
                            decoration: InputDecoration(
                              labelText: l10n.defenseSummary,
                              alignLabelWithHint: true,
                            ),
                          ),
                          Align(
                            alignment: Alignment.centerLeft,
                            child: TextButton.icon(
                              key: const ValueKey('draft-protocol'),
                              onPressed: _busy ? null : _draftProtocol,
                              icon: const Icon(Icons.auto_awesome_outlined),
                              label: Text(l10n.defenseDraftProtocol),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: Spacing.xl),
                    Wrap(
                      spacing: Spacing.s,
                      runSpacing: Spacing.s,
                      children: [
                        OutlinedButton(
                          key: const ValueKey('save-defense'),
                          onPressed: _busy ? null : _save,
                          child: Text(l10n.actionSave),
                        ),
                        FilledButton(
                          key: const ValueKey('finish-defense'),
                          onPressed: _busy ? null : _finish,
                          child: Text(l10n.defenseFinish),
                        ),
                      ],
                    ),
                  ],
                ),
        ),
      ),
    );
  }
}
