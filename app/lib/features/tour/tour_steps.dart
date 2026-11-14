import 'package:material_ui/material_ui.dart';

import '../../l10n/generated/app_localizations.dart';
import 'tour_step.dart';

class TourSteps {
  const TourSteps._();

  static List<TourStep> of(AppLocalizations l10n) => [
    TourStep(
      icon: Icons.wb_sunny_outlined,
      title: l10n.tourTodayTitle,
      body: l10n.tourTodayBody,
    ),
    TourStep(
      icon: Icons.explore_outlined,
      title: l10n.tourNavigatorTitle,
      body: l10n.tourNavigatorBody,
    ),
    TourStep(
      icon: Icons.hub_outlined,
      title: l10n.tourNotesTitle,
      body: l10n.tourNotesBody,
    ),
    TourStep(
      icon: Icons.forum_outlined,
      title: l10n.tourChatsTitle,
      body: l10n.tourChatsBody,
    ),
    TourStep(
      icon: Icons.palette_outlined,
      title: l10n.tourStyleTitle,
      body: l10n.tourStyleBody,
    ),
  ];
}
