import '../../l10n/generated/app_localizations.dart';
import 'source_kind.dart';

class SourceLabels {
  const SourceLabels(this._l10n);

  final AppLocalizations _l10n;

  String kind(SourceKind kind) => switch (kind) {
    SourceKind.book => _l10n.sourceKindBook,
    SourceKind.article => _l10n.sourceKindArticle,
    SourceKind.chapter => _l10n.sourceKindChapter,
    SourceKind.web => _l10n.sourceKindWeb,
    SourceKind.thesis => _l10n.sourceKindThesis,
    SourceKind.other => _l10n.sourceKindOther,
  };
}
