import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/offline/client_id.dart';
import '../../design/components/form_error.dart';
import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'brief.dart';
import 'brief_banner.dart';
import 'dialog_scope.dart';
import 'intent_sheet.dart';
import 'letter_sheet.dart';
import 'message_bubble.dart';
import 'message_composer.dart';
import 'navigator_message.dart';
import 'navigator_providers.dart';
import 'pending_bubble.dart';
import 'question.dart';
import 'reply_actions.dart';
import 'teacher_letter.dart';

class DialogScreen extends ConsumerStatefulWidget {
  const DialogScreen({this.dialogId, this.scope, super.key});

  final String? dialogId;
  final DialogScope? scope;

  @override
  ConsumerState<DialogScreen> createState() => _DialogScreenState();
}

class _DialogScreenState extends ConsumerState<DialogScreen> {
  late final String _id = widget.dialogId ?? ClientId.generate();
  late bool _created = widget.dialogId != null;
  late bool _loading = widget.dialogId != null;
  final _scroll = ScrollController();
  final _messages = <NavigatorMessage>[];
  Brief? _brief;
  String _title = '';
  Object? _loadError;
  Question? _pending;
  String? _error;

  @override
  void initState() {
    super.initState();
    if (_created) {
      _load();
    }
  }

  @override
  void dispose() {
    _scroll.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final loaded = await ref.read(navigatorApiProvider).dialog(_id);
      if (!mounted) {
        return;
      }
      setState(() {
        _messages
          ..clear()
          ..addAll(loaded.value.messages);
        _brief = loaded.value.brief;
        _title = loaded.value.dialog.title;
        _loading = false;
      });
      _toBottom();
    } on Object catch (error) {
      if (mounted) {
        setState(() {
          _loadError = error;
          _loading = false;
        });
      }
    }
  }

  Future<void> _ask(Question question) async {
    final l10n = context.l10n;
    final api = ref.read(navigatorApiProvider);
    setState(() {
      _pending = question;
      _error = null;
    });
    _toBottom();
    try {
      if (!_created) {
        final details = await api.create(_id, widget.scope);
        _created = true;
        _brief = details.brief;
      }
      final answer = await api.ask(_id, question);
      if (!mounted) {
        return;
      }
      setState(() {
        _messages
          ..add(answer.question)
          ..add(answer.reply);
        if (_title.isEmpty) {
          _title = answer.question.text;
        }
        _pending = null;
      });
      ref.invalidate(dialogsProvider);
      _toBottom();
    } on Object catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    }
  }

  void _send(String text) =>
      _ask(Question(key: ClientId.generate(), text: text));

  void _more() =>
      _ask(Question(key: ClientId.generate(), text: '', more: true));

  Future<void> _misread() async {
    final asked = _messages.lastWhere(
      (m) => !m.fromNavigator && m.text.isNotEmpty,
      orElse: () => _messages.first,
    );
    final intent = await IntentSheet.show(context);
    if (intent != null && mounted) {
      await _ask(
        Question(key: ClientId.generate(), text: asked.text, intent: intent),
      );
    }
  }

  Future<void> _delete() async {
    final l10n = context.l10n;
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(l10n.navigatorDeleteConfirm),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: Text(l10n.actionCancel),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: Text(l10n.navigatorDelete),
          ),
        ],
      ),
    );
    if (!(confirmed ?? false) || !mounted) {
      return;
    }
    try {
      await ref.read(navigatorApiProvider).delete(_id);
      ref.invalidate(dialogsProvider);
      if (mounted) {
        context.pop();
      }
    } on Object catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    }
  }

  void _letter() {
    LetterSheet.show(
      context,
      TeacherLetter.compose(context.l10n, _brief, _messages),
    );
  }

  void _toBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scroll.hasClients) {
        _scroll.jumpTo(_scroll.position.maxScrollExtent);
      }
    });
  }

  List<Widget> _items(BuildContext context) {
    final pending = _pending;
    final brief = _brief;
    final last = _messages.isEmpty ? null : _messages.last;
    return [
      if (brief != null) ...[
        BriefBanner(brief: brief),
        const SizedBox(height: Spacing.l),
      ],
      if (_messages.isEmpty && pending == null)
        Text(
          context.l10n.navigatorEmptyDialog,
          style: TextStyle(color: context.palette.muted),
        ),
      for (final message in _messages) ...[
        MessageBubble(key: ValueKey(message.seq), message: message),
        const SizedBox(height: Spacing.m),
      ],
      if (pending != null)
        PendingBubble(
          text: pending.text,
          error: _error,
          onRetry: () => _ask(pending),
        )
      else if (last != null && last.fromNavigator)
        ReplyActions(reply: last, onMore: _more, onMisread: _misread),
      if (pending == null && _error != null) FormError(message: _error!),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final loadError = _loadError;
    return Scaffold(
      appBar: AppBar(
        title: Text(
          _title.isEmpty ? l10n.navigatorNewQuestion : _title,
          overflow: TextOverflow.ellipsis,
        ),
        actions: [
          if (_brief != null && _messages.isNotEmpty)
            IconButton(
              tooltip: l10n.letterTitle,
              icon: const Icon(Icons.mail_outline),
              onPressed: _letter,
            ),
          if (_created)
            IconButton(
              tooltip: l10n.navigatorDelete,
              icon: const Icon(Icons.delete_outline),
              onPressed: _pending == null ? _delete : null,
            ),
        ],
      ),
      body: SafeArea(
        top: false,
        child: Align(
          alignment: Alignment.topCenter,
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 760),
            child: Column(
              children: [
                Expanded(
                  child: _loading
                      ? const Center(child: CircularProgressIndicator())
                      : loadError != null
                      ? Padding(
                          padding: const EdgeInsets.all(Spacing.l),
                          child: FormError(
                            message: ApiErrorText.of(l10n, loadError),
                          ),
                        )
                      : ListView(
                          controller: _scroll,
                          padding: const EdgeInsets.all(Spacing.l),
                          children: _items(context),
                        ),
                ),
                MessageComposer(
                  enabled: !_loading && loadError == null && _pending == null,
                  onSend: _send,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
