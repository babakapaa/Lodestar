import 'source_kind.dart';

class Source {
  const Source({
    required this.id,
    required this.kind,
    required this.title,
    this.authors = const [],
    this.year = '',
    this.container = '',
    this.publisher = '',
    this.city = '',
    this.volume = '',
    this.issue = '',
    this.pages = '',
    this.url = '',
    this.doi = '',
    this.isbn = '',
    this.accessedAt = '',
    this.note = '',
  });

  factory Source.fromJson(Map<String, Object?> json, {String? id}) => Source(
    id: id ?? json['id']! as String,
    kind: SourceKind.parse(json['kind']! as String),
    title: json['title']! as String,
    authors: [
      for (final author in json['authors']! as List<Object?>) author! as String,
    ],
    year: json['year']! as String,
    container: json['container']! as String,
    publisher: json['publisher']! as String,
    city: json['city']! as String,
    volume: json['volume']! as String,
    issue: json['issue']! as String,
    pages: json['pages']! as String,
    url: json['url']! as String,
    doi: json['doi']! as String,
    isbn: json['isbn']! as String,
    accessedAt: json['accessedAt']! as String,
    note: json['note']! as String,
  );

  final String id;
  final SourceKind kind;
  final String title;
  final List<String> authors;
  final String year;
  final String container;
  final String publisher;
  final String city;
  final String volume;
  final String issue;
  final String pages;
  final String url;
  final String doi;
  final String isbn;
  final String accessedAt;
  final String note;

  Map<String, Object?> toJson() => {
    'id': id,
    'kind': kind.name,
    'title': title,
    'authors': authors,
    'year': year,
    'container': container,
    'publisher': publisher,
    'city': city,
    'volume': volume,
    'issue': issue,
    'pages': pages,
    'url': url,
    'doi': doi,
    'isbn': isbn,
    'accessedAt': accessedAt,
    'note': note,
  };
}
