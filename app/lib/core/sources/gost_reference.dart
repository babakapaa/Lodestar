import 'source.dart';
import 'source_kind.dart';

abstract final class GostReference {
  static const _dash = ' – ';
  static final _cyrillic = RegExp('[а-яёәғқңөұүһі]', caseSensitive: false);

  static String format(Source source) {
    final cyrillic = _cyrillic.hasMatch(source.title);
    final authors = source.authors;
    final head = authors.isNotEmpty && authors.length <= 3
        ? '${_heading(authors.first)} '
        : '';
    final responsibility = switch (authors.length) {
      0 => '',
      > 4 =>
        ' / ${authors.take(3).map(_initialsFirst).join(', ')} '
            '${cyrillic ? '[и др.]' : '[et al.]'}',
      _ => ' / ${authors.map(_initialsFirst).join(', ')}',
    };
    final title = switch (source.kind) {
      SourceKind.thesis =>
        '${source.title} : ${cyrillic ? 'диссертация' : 'thesis'}',
      SourceKind.web =>
        '${source.title} : ${cyrillic ? 'электронный ресурс' : 'electronic resource'}',
      _ => source.title,
    };
    final parts = <String>['$head$title$responsibility'];
    final container = source.container.trim();
    switch (source.kind) {
      case SourceKind.article || SourceKind.chapter || SourceKind.web
          when container.isNotEmpty:
        parts[0] = '${parts[0]} // $container';
        parts.addAll([
          if (source.year.isNotEmpty) source.year,
          ?_volumeIssue(source, cyrillic),
          if (source.pages.isNotEmpty)
            '${cyrillic ? 'С.' : 'P.'} ${source.pages}',
        ]);
      default:
        final place = [
          if (source.city.isNotEmpty) source.city,
          if (source.publisher.isNotEmpty) source.publisher,
        ].join(' : ');
        final imprint = [
          if (place.isNotEmpty) place,
          if (source.year.isNotEmpty) source.year,
        ].join(', ');
        parts.addAll([
          if (imprint.isNotEmpty) imprint,
          ?_volumeIssue(source, cyrillic),
          if (source.pages.isNotEmpty)
            '${cyrillic ? 'С.' : 'P.'} ${source.pages}',
        ]);
    }
    if (source.isbn.isNotEmpty) {
      parts.add('ISBN ${source.isbn}');
    }
    if (source.doi.isNotEmpty) {
      parts.add('DOI ${source.doi}');
    }
    if (source.url.isNotEmpty && source.doi.isEmpty) {
      final accessed = _date(source.accessedAt);
      parts.add(
        'URL: ${source.url}'
        '${accessed == null ? '' : ' (${cyrillic ? 'дата обращения' : 'accessed'}: $accessed)'}',
      );
    }
    return '${parts.map(_trimDot).join('.$_dash')}.';
  }

  static String list(List<Source> sources) {
    final sorted = [...sources]
      ..sort(
        (a, b) =>
            _sortKey(a).toLowerCase().compareTo(_sortKey(b).toLowerCase()),
      );
    return [
      for (var i = 0; i < sorted.length; i++) '${i + 1}. ${format(sorted[i])}',
    ].join('\n');
  }

  static String _sortKey(Source s) =>
      s.authors.isNotEmpty && s.authors.length <= 3 ? s.authors.first : s.title;

  static String? _volumeIssue(Source s, bool cyrillic) {
    final volume = s.volume.isEmpty
        ? null
        : '${cyrillic ? 'Т.' : 'Vol.'} ${s.volume}';
    final issue = s.issue.isEmpty
        ? null
        : '${cyrillic ? '№' : 'No.'} ${s.issue}';
    final both = [?volume, ?issue].join(', ');
    return both.isEmpty ? null : both;
  }

  static String _heading(String author) {
    final (family, initials) = _split(author);
    return initials.isEmpty ? family : '$family, $initials';
  }

  static String _initialsFirst(String author) {
    final (family, initials) = _split(author);
    return initials.isEmpty ? family : '$initials $family';
  }

  static (String, String) _split(String author) {
    final at = author.indexOf(',');
    if (at < 0) {
      return (author.trim(), '');
    }
    final family = author.substring(0, at).trim();
    final given = author.substring(at + 1).trim();
    final initials = [
      for (final part in given.split(RegExp(r'[\s.]+')))
        if (part.isNotEmpty)
          part.contains('-')
              ? part.split('-').map((p) => '${p[0]}.').join('-')
              : '${part[0]}.',
    ].join(' ');
    return (family, initials);
  }

  static String? _date(String iso) {
    final date = DateTime.tryParse(iso);
    if (date == null) {
      return null;
    }
    String two(int n) => n.toString().padLeft(2, '0');
    return '${two(date.day)}.${two(date.month)}.${date.year}';
  }

  static String _trimDot(String part) =>
      part.endsWith('.') ? part.substring(0, part.length - 1) : part;
}
