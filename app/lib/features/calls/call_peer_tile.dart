import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import '../../design/tokens/tokens_context.dart';
import '../../l10n/l10n_context.dart';
import 'call_peer.dart';

class CallPeerTile extends StatelessWidget {
  const CallPeerTile({required this.peer, required this.video, super.key});

  final CallPeer peer;
  final Widget video;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final palette = tokens.palette;
    final name = peer.local ? context.l10n.callYou(peer.name) : peer.name;
    return Semantics(
      label: name,
      child: DecoratedBox(
        decoration: BoxDecoration(
          color: palette.surface2,
          borderRadius: tokens.borderRadius,
          border: Border.all(
            color: peer.speaking ? palette.accent : palette.border,
            width: peer.speaking ? 3 : 1,
          ),
        ),
        child: ClipRRect(
          borderRadius: tokens.borderRadius,
          child: Stack(
            fit: StackFit.expand,
            children: [
              if (peer.video)
                video
              else
                Center(
                  child: CircleAvatar(
                    radius: 32,
                    backgroundColor: palette.accentSoft,
                    child: Text(
                      peer.name.isEmpty ? '?' : peer.name.characters.first,
                      style: Theme.of(context).textTheme.headlineSmall,
                    ),
                  ),
                ),
              Positioned(
                left: Spacing.s,
                right: Spacing.s,
                bottom: Spacing.s,
                child: Row(
                  children: [
                    if (!peer.microphone)
                      Icon(Icons.mic_off, size: 16, color: palette.warn),
                    if (peer.screen)
                      Icon(Icons.screen_share, size: 16, color: palette.accent),
                    const SizedBox(width: Spacing.xs),
                    Expanded(
                      child: Text(
                        name,
                        overflow: TextOverflow.ellipsis,
                        style: Theme.of(context).textTheme.labelMedium,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
