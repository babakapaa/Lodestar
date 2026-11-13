import '../../l10n/generated/app_localizations.dart';
import 'brief.dart';
import 'navigator_message.dart';

abstract final class TeacherLetter {
  static String compose(
    AppLocalizations l10n,
    Brief? brief,
    List<NavigatorMessage> messages,
  ) {
    final questions = [
      for (final m in messages)
        if (!m.fromNavigator && m.text.isNotEmpty) m.text,
    ];
    final sources = {
      for (final m in messages)
        for (final c in m.citations)
          c.reference.isEmpty ? c.title : c.reference,
    };
    return [
      l10n.letterGreeting,
      if (brief != null)
        brief.isCourse
            ? l10n.letterCourse(brief.title)
            : l10n.letterAssignment(brief.title, brief.course),
      if (questions.isNotEmpty) '${l10n.letterQuestion}\n${questions.first}',
      if (questions.length > 1)
        [
          l10n.letterTried,
          for (final q in questions.skip(1)) '- $q',
        ].join('\n'),
      if (sources.isNotEmpty)
        [l10n.letterSources, for (final s in sources) '- $s'].join('\n'),
      l10n.letterStuck,
      l10n.letterThanks,
    ].join('\n\n');
  }
}
