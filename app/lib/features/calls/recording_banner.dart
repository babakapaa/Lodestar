import 'package:material_ui/material_ui.dart';

import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';

class RecordingBanner extends StatelessWidget {
  const RecordingBanner({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      key: const ValueKey('call-recording'),
      padding: const EdgeInsets.fromLTRB(Spacing.m, Spacing.m, Spacing.m, 0),
      child: NoticeBanner(
        icon: Icons.fiber_manual_record,
        message: context.l10n.callRecording,
      ),
    );
  }
}
