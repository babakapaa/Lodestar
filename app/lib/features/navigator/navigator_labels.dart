import '../../l10n/generated/app_localizations.dart';
import 'hint_step.dart';
import 'navigator_intent.dart';

class NavigatorLabels {
  const NavigatorLabels(this._l10n);

  final AppLocalizations _l10n;

  String intent(NavigatorIntent intent) => switch (intent) {
    NavigatorIntent.understand => _l10n.intentUnderstand,
    NavigatorIntent.find => _l10n.intentFind,
    NavigatorIntent.check => _l10n.intentCheck,
    NavigatorIntent.doForMe => _l10n.intentDoForMe,
    NavigatorIntent.organize => _l10n.intentOrganize,
    NavigatorIntent.state || NavigatorIntent.closed => _l10n.intentState,
  };

  String step(HintStep step) =>
      _l10n.navigatorStep(step.index + 1, switch (step) {
        HintStep.question => _l10n.stepQuestion,
        HintStep.direction => _l10n.stepDirection,
        HintStep.source => _l10n.stepSource,
        HintStep.fragment => _l10n.stepFragment,
        HintStep.example => _l10n.stepExample,
      });
}
