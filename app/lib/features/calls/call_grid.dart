import 'dart:math';

import 'package:material_ui/material_ui.dart';

import '../../design/tokens/spacing.dart';
import 'call_engine.dart';
import 'call_peer_tile.dart';

class CallGrid extends StatelessWidget {
  const CallGrid({required this.engine, super.key});

  final CallEngine engine;

  @override
  Widget build(BuildContext context) {
    final peers = engine.peers;
    return LayoutBuilder(
      builder: (context, constraints) {
        final columns = max(
          1,
          min(sqrt(peers.length).ceil(), (constraints.maxWidth / 220).floor()),
        );
        return GridView.builder(
          padding: const EdgeInsets.all(Spacing.m),
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: columns,
            mainAxisSpacing: Spacing.m,
            crossAxisSpacing: Spacing.m,
            childAspectRatio: 4 / 3,
          ),
          itemCount: peers.length,
          itemBuilder: (context, index) => CallPeerTile(
            key: ValueKey('call-peer-${peers[index].id}'),
            peer: peers[index],
            video: engine.view(peers[index]),
          ),
        );
      },
    );
  }
}
