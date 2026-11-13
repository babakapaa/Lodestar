import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'call_engine.dart';

class CallControls extends StatelessWidget {
  const CallControls({
    required this.engine,
    required this.screenShare,
    required this.onHangUp,
    super.key,
  });

  final CallEngine engine;
  final bool screenShare;
  final VoidCallback onHangUp;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final palette = context.tokens.palette;
    return Wrap(
      alignment: WrapAlignment.center,
      spacing: Spacing.m,
      children: [
        IconButton.filledTonal(
          key: const ValueKey('call-microphone'),
          tooltip: engine.microphone ? l10n.callMute : l10n.callUnmute,
          icon: Icon(engine.microphone ? Icons.mic : Icons.mic_off),
          onPressed: () => engine.setMicrophone(on: !engine.microphone),
        ),
        IconButton.filledTonal(
          key: const ValueKey('call-camera'),
          tooltip: engine.camera ? l10n.callCameraOff : l10n.callCameraOn,
          icon: Icon(engine.camera ? Icons.videocam : Icons.videocam_off),
          onPressed: () => engine.setCamera(on: !engine.camera),
        ),
        if (screenShare)
          IconButton.filledTonal(
            key: const ValueKey('call-screen'),
            tooltip: engine.screen ? l10n.callScreenOff : l10n.callScreenOn,
            icon: Icon(
              engine.screen ? Icons.stop_screen_share : Icons.screen_share,
            ),
            onPressed: () => engine.setScreen(on: !engine.screen),
          ),
        IconButton.filled(
          key: const ValueKey('call-hang-up'),
          tooltip: l10n.callHangUp,
          style: IconButton.styleFrom(
            backgroundColor: palette.warn,
            foregroundColor: palette.onAccent,
          ),
          icon: const Icon(Icons.call_end),
          onPressed: onHangUp,
        ),
      ],
    );
  }
}
