import 'dart:async';

import 'package:livekit_client/livekit_client.dart';
import 'package:material_ui/material_ui.dart';

import 'call_engine.dart';
import 'call_peer.dart';

class LiveKitCallEngine extends CallEngine {
  LiveKitCallEngine() {
    _listener.listen((_) => notifyListeners());
  }

  final Room _room = Room(
    roomOptions: const RoomOptions(adaptiveStream: true, dynacast: true),
  );
  late final EventsListener<RoomEvent> _listener = _room.createListener();

  @override
  List<CallPeer> get peers => [
    if (_room.localParticipant case final local?) _peer(local, local: true),
    for (final remote in _room.remoteParticipants.values)
      _peer(remote, local: false),
  ];

  @override
  bool get microphone => _room.localParticipant?.isMicrophoneEnabled() ?? false;

  @override
  bool get camera => _room.localParticipant?.isCameraEnabled() ?? false;

  @override
  bool get screen => _room.localParticipant?.isScreenShareEnabled() ?? false;

  @override
  Future<void> connect(String url, String token) async {
    await _room.connect(url, token);
    notifyListeners();
  }

  @override
  Future<void> setMicrophone({required bool on}) async {
    await _room.localParticipant?.setMicrophoneEnabled(on);
    notifyListeners();
  }

  @override
  Future<void> setCamera({required bool on}) async {
    await _room.localParticipant?.setCameraEnabled(on);
    notifyListeners();
  }

  @override
  Future<void> setScreen({required bool on}) async {
    await _room.localParticipant?.setScreenShareEnabled(on);
    notifyListeners();
  }

  @override
  Future<void> leave() => _room.disconnect();

  @override
  Widget view(CallPeer peer) {
    final Participant? participant = peer.local
        ? _room.localParticipant
        : _room.remoteParticipants[peer.id];
    final track = participant == null ? null : _track(participant);
    if (track == null) {
      return const SizedBox.shrink();
    }
    return VideoTrackRenderer(track, fit: VideoViewFit.contain);
  }

  @override
  void dispose() {
    unawaited(_listener.dispose());
    unawaited(_room.dispose());
    super.dispose();
  }

  static VideoTrack? _track(Participant participant) {
    final publications = participant.videoTrackPublications
        .where((p) => !p.muted && p.track is VideoTrack)
        .toList();
    if (publications.isEmpty) {
      return null;
    }
    final shown = publications.firstWhere(
      (p) => p.isScreenShare,
      orElse: () => publications.first,
    );
    return shown.track! as VideoTrack;
  }

  static CallPeer _peer(Participant participant, {required bool local}) =>
      CallPeer(
        id: participant.identity,
        name: participant.name.isEmpty
            ? participant.identity
            : participant.name,
        local: local,
        speaking: participant.isSpeaking,
        microphone: participant.isMicrophoneEnabled(),
        video: _track(participant) != null,
        screen: participant.isScreenShareEnabled(),
      );
}
