import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:material_ui/material_ui.dart';

import '../../core/api/api_error_text.dart';
import '../../core/api/api_exception.dart';
import '../../design/components/empty_state.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'call_controls.dart';
import 'call_engine.dart';
import 'call_grid.dart';
import 'call_ticket.dart';
import 'calls_api.dart';
import 'calls_providers.dart';
import 'recording_banner.dart';

class CallScreen extends ConsumerStatefulWidget {
  const CallScreen({required this.chatId, super.key});

  final String chatId;

  @override
  ConsumerState<CallScreen> createState() => _CallScreenState();
}

class _CallScreenState extends ConsumerState<CallScreen> {
  late final CallEngine _engine;
  late final CallsApi _api;
  CallTicket? _ticket;
  String? _error;
  bool _recording = false;
  bool _left = false;

  @override
  void initState() {
    super.initState();
    _engine = ref.read(callEngineFactoryProvider)();
    _api = ref.read(callsApiProvider);
    WidgetsBinding.instance.addPostFrameCallback((_) => unawaited(_join()));
  }

  @override
  void dispose() {
    unawaited(_leave().whenComplete(_engine.dispose));
    super.dispose();
  }

  Future<void> _join() async {
    final l10n = context.l10n;
    try {
      final ticket = await _api.join(widget.chatId);
      if (!mounted) {
        return;
      }
      setState(() {
        _ticket = ticket;
        _recording = ticket.call.recording;
      });
      await _engine.connect(ticket.url, ticket.token);
      await _engine.setMicrophone(on: true);
    } on ApiException catch (error) {
      if (mounted) {
        setState(() => _error = ApiErrorText.of(l10n, error));
      }
    }
  }

  Future<void> _leave() async {
    if (_left) {
      return;
    }
    _left = true;
    final ticket = _ticket;
    await _engine.leave();
    if (ticket != null) {
      try {
        await _api.leave(ticket.call.id);
      } on ApiException {
        return;
      }
    }
  }

  Future<void> _record() async {
    final ticket = _ticket;
    if (ticket == null) {
      return;
    }
    final l10n = context.l10n;
    final messenger = ScaffoldMessenger.of(context);
    try {
      final call = _recording
          ? await _api.stopRecording(ticket.call.id)
          : await _api.startRecording(ticket.call.id);
      if (mounted) {
        setState(() => _recording = call.recording);
      }
    } on ApiException catch (error) {
      messenger.showSnackBar(
        SnackBar(content: Text(ApiErrorText.of(l10n, error))),
      );
    }
  }

  Future<void> _hangUp() async {
    await _leave();
    if (mounted) {
      context.pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final ticket = _ticket;
    final error = _error;
    ref.listen(activeCallProvider(widget.chatId), (_, next) {
      final call = next.value;
      if (call != null && call.id == ticket?.call.id) {
        setState(() => _recording = call.recording);
      }
    });
    return Scaffold(
      appBar: AppBar(title: Text(l10n.callTitle)),
      body: SafeArea(
        top: false,
        child: error != null
            ? EmptyState(
                key: const ValueKey('call-error'),
                icon: Icons.phone_disabled_outlined,
                title: l10n.callFailed,
                body: error,
              )
            : ticket == null
            ? const Center(child: CircularProgressIndicator())
            : ListenableBuilder(
                listenable: _engine,
                builder: (context, _) => Column(
                  children: [
                    if (_recording) const RecordingBanner(),
                    Expanded(child: CallGrid(engine: _engine)),
                    Padding(
                      padding: const EdgeInsets.all(Spacing.l),
                      child: CallControls(
                        engine: _engine,
                        screenShare: ticket.screenShare,
                        record: ticket.record,
                        recording: _recording,
                        onRecord: _record,
                        onHangUp: _hangUp,
                      ),
                    ),
                  ],
                ),
              ),
      ),
    );
  }
}
