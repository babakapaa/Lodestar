import 'package:lodestar/features/calls/call_engine.dart';
import 'package:lodestar/features/calls/call_peer.dart';
import 'package:material_ui/material_ui.dart';

class FakeCallEngine extends CallEngine {
  String? url;
  String? token;
  bool left = false;
  final List<CallPeer> others = [];

  @override
  bool microphone = false;

  @override
  bool camera = false;

  @override
  bool screen = false;

  @override
  List<CallPeer> get peers => [
    if (url != null)
      CallPeer(
        id: 'u1',
        name: 'Айгерим Сапарова',
        local: true,
        speaking: false,
        microphone: microphone,
        video: camera,
        screen: screen,
      ),
    ...others,
  ];

  void join(CallPeer peer) {
    others.add(peer);
    notifyListeners();
  }

  @override
  Future<void> connect(String url, String token) async {
    this.url = url;
    this.token = token;
    notifyListeners();
  }

  @override
  Future<void> setMicrophone({required bool on}) async {
    microphone = on;
    notifyListeners();
  }

  @override
  Future<void> setCamera({required bool on}) async {
    camera = on;
    notifyListeners();
  }

  @override
  Future<void> setScreen({required bool on}) async {
    screen = on;
    notifyListeners();
  }

  @override
  Future<void> leave() async => left = true;

  @override
  Widget view(CallPeer peer) => ColoredBox(
    key: ValueKey('call-video-${peer.id}'),
    color: const Color(0xFF000000),
  );
}
