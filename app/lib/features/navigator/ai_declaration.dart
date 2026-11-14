import 'package:intl/intl.dart';

import '../../l10n/generated/app_localizations.dart';
import 'brief.dart';
import 'navigator_labels.dart';
import 'navigator_message.dart';

abstract final class AiDeclaration {
  static String compose(
    AppLocalizations l10n,
    NavigatorLabels labels,
    String locale,
    Brief? brief,
    List<NavigatorMessage> messages,
  ) {
    final questions = [
      for (final m in messages)
        if (!m.fromNavigator && m.text.isNotEmpty) m.text,
    ];
    final help = {
      for (final m in messages)
        if (m.fromNavigator && m.step != null) labels.step(m.step!),
    };
    final sources = {
      for (final m in messages)
        for (final c in m.citations)
          c.reference.isEmpty ? c.title : c.reference,
    };
    String section(String title, Iterable<String> lines) =>
        ['## $title', for (final line in lines) '- $line'].join('\n');
    return [
      '# ${l10n.declarationTitle}',
      l10n.declarationIntro,
      if (brief != null)
        '**${l10n.declarationWork}:** ${brief.isCourse ? brief.title : '${brief.title}, ${brief.course}'}',
      if (questions.isNotEmpty) section(l10n.declarationQuestions, questions),
      if (help.isNotEmpty) section(l10n.declarationHelp, help),
      if (sources.isNotEmpty) section(l10n.declarationSources, sources),
      '## ${l10n.declarationOwn}\n${l10n.declarationOwnHint}',
      '**${l10n.declarationDate}:** ${DateFormat.yMMMMd(locale).format(DateTime.now())}',
    ].join('\n\n');
  }
}
