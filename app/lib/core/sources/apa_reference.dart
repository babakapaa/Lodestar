import 'author_name.dart';
import 'source.dart';
import 'source_kind.dart';

abstract final class ApaReference {
  static const _maxAuthors = 20;

  static String format(Source source) {
    final year = '(${source.year.isEmpty ? 'n.d.' : source.year}).';
    final title = '${_title(source)}.';
    final authors = _authors(source.authors);
    return [
      if (authors.isEmpty) ...[title, year] else ...[authors, year, title],
      ?_container(source),
      ?_link(source),
    ].join(' ');
  }

  static String list(List<Source> sources) {
    final sorted = [...sources]
      ..sort(
        (a, b) =>
            _sortKey(a).toLowerCase().compareTo(_sortKey(b).toLowerCase()),
      );
    return sorted.map(format).join('\n');
  }

  static String _sortKey(Source s) =>
      s.authors.isNotEmpty ? AuthorName.parse(s.authors.first).family : s.title;

  static String _title(Source source) {
    final title = _trimDot(source.title.trim());
    if (source.kind != SourceKind.thesis) {
      return title;
    }
    final place = source.publisher.trim();
    return '$title [Thesis${place.isEmpty ? '' : ', $place'}]';
  }

  static String _authors(List<String> authors) {
    final names = [for (final author in authors) _name(author)];
    if (names.isEmpty) {
      return '';
    }
    final String joined;
    if (names.length == 1) {
      joined = names.first;
    } else if (names.length > _maxAuthors) {
      joined = '${names.take(_maxAuthors - 1).join(', ')}, ... ${names.last}';
    } else {
      joined = '${names.take(names.length - 1).join(', ')}, & ${names.last}';
    }
    return joined.endsWith('.') ? joined : '$joined.';
  }

  static String _name(String author) {
    final name = AuthorName.parse(author);
    return name.initials.isEmpty
        ? name.family
        : '${name.family}, ${name.initials}';
  }

  static String? _container(Source source) {
    final container = _trimDot(source.container.trim());
    final publisher = _trimDot(source.publisher.trim());
    final pages = source.pages.trim();
    switch (source.kind) {
      case SourceKind.article when container.isNotEmpty:
        final volume = source.volume.trim();
        final issue = source.issue.trim();
        final parts = [
          container,
          if (volume.isNotEmpty || issue.isNotEmpty)
            '$volume${issue.isEmpty ? '' : '($issue)'}',
          if (pages.isNotEmpty) pages,
        ];
        return '${parts.join(', ')}.';
      case SourceKind.chapter when container.isNotEmpty:
        final chapter = 'In $container${pages.isEmpty ? '' : ' (pp. $pages)'}.';
        return publisher.isEmpty ? chapter : '$chapter $publisher.';
      case SourceKind.web when container.isNotEmpty:
        return '$container.';
      case SourceKind.thesis:
        return null;
      default:
        if (publisher.isNotEmpty) {
          return '$publisher.';
        }
        return container.isEmpty ? null : '$container.';
    }
  }

  static String? _link(Source source) {
    if (source.doi.isNotEmpty) {
      return 'https://doi.org/${source.doi}';
    }
    return source.url.isEmpty ? null : source.url;
  }

  static String _trimDot(String part) =>
      part.endsWith('.') ? part.substring(0, part.length - 1) : part;
}
