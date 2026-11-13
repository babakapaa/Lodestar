import 'package:material_ui/material_ui.dart';

import 'call_peer.dart';

abstract class CallEngine extends ChangeNotifier {
  List<CallPeer> get peers;

  bool get microphone;

  bool get camera;

  bool get screen;

  Future<void> connect(String url, String token);

  Future<void> setMicrophone({required bool on});

  Future<void> setCamera({required bool on});

  Future<void> setScreen({required bool on});

  Future<void> leave();

  Widget view(CallPeer peer);
}
