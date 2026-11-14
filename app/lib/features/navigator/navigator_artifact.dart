import 'package:flutter/foundation.dart';

@immutable
class NavigatorArtifact {
  const NavigatorArtifact({required this.title, required this.text});

  final String title;
  final String text;
}
