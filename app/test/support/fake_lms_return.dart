import 'package:lodestar/features/lti/lms_return.dart';

class FakeLmsReturn implements LmsReturn {
  final List<Uri> opened = [];

  @override
  Future<void> open(Uri url) async => opened.add(url);
}
