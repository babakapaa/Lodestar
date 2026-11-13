import 'package:material_ui/material_ui.dart';

import '../../design/components/notice_banner.dart';
import '../../design/tokens/spacing.dart';
import '../../l10n/l10n_context.dart';
import 'call_consent.dart';

class ConsentBanner extends StatelessWidget {
  const ConsentBanner({
    required this.consent,
    required this.me,
    required this.onAnswer,
    super.key,
  });

  final CallConsent consent;
  final String? me;
  final ValueChanged<bool> onAnswer;

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final asked = consent.waiting.contains(me);
    return Padding(
      key: const ValueKey('call-consent'),
      padding: const EdgeInsets.fromLTRB(Spacing.m, Spacing.m, Spacing.m, 0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          NoticeBanner(
            icon: Icons.how_to_vote_outlined,
            message: asked
                ? (consent.audioOnly
                      ? l10n.callConsentAskAudio
                      : l10n.callConsentAsk)
                : l10n.callConsentWaiting(consent.waiting.length),
          ),
          if (asked)
            Wrap(
              spacing: Spacing.s,
              children: [
                FilledButton(
                  key: const ValueKey('consent-agree'),
                  onPressed: () => onAnswer(true),
                  child: Text(l10n.callConsentAgree),
                ),
                OutlinedButton(
                  key: const ValueKey('consent-decline'),
                  onPressed: () => onAnswer(false),
                  child: Text(l10n.callConsentDecline),
                ),
              ],
            ),
        ],
      ),
    );
  }
}
