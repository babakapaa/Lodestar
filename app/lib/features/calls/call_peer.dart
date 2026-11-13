class CallPeer {
  const CallPeer({
    required this.id,
    required this.name,
    required this.local,
    required this.speaking,
    required this.microphone,
    required this.video,
    required this.screen,
  });

  final String id;
  final String name;
  final bool local;
  final bool speaking;
  final bool microphone;
  final bool video;
  final bool screen;
}
