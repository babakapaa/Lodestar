import 'package:url_launcher/url_launcher.dart';

import 'lms_return.dart';

class BrowserLmsReturn implements LmsReturn {
  const BrowserLmsReturn();

  @override
  Future<void> open(Uri url) async {
    await launchUrl(url, webOnlyWindowName: '_self');
  }
}
