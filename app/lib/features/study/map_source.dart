import '../../core/sources/gost_reference.dart';
import '../../core/sources/source.dart';
import '../../core/sources/source_kind.dart';

class MapSource {
  const MapSource({
    required this.sourceId,
    required this.reference,
    this.kind = '',
    this.doi = '',
    this.url = '',
    this.found = '',
    this.why = '',
    this.taken = '',
  });

  factory MapSource.fromSource(Source source) => MapSource(
    sourceId: source.id,
    reference: GostReference.format(source),
    kind: source.kind.name,
    doi: source.doi,
    url: source.url,
    found: _found(source),
  );

  factory MapSource.fromJson(Map<String, Object?> json) => MapSource(
    sourceId: json['sourceId'] as String? ?? '',
    reference: json['reference'] as String? ?? '',
    kind: json['kind'] as String? ?? '',
    doi: json['doi'] as String? ?? '',
    url: json['url'] as String? ?? '',
    found: json['found'] as String? ?? '',
    why: json['why'] as String? ?? '',
    taken: json['taken'] as String? ?? '',
  );

  static List<MapSource> list(Object? json) => [
    for (final item in json as List<Object?>? ?? const [])
      MapSource.fromJson(item! as Map<String, Object?>),
  ];

  static const _scholarlyKinds = {
    SourceKind.article,
    SourceKind.chapter,
    SourceKind.thesis,
  };

  static String _found(Source source) {
    if (source.doi.isNotEmpty) {
      return 'doi.org';
    }
    return Uri.tryParse(source.url)?.host ?? '';
  }

  final String sourceId;
  final String reference;
  final String kind;
  final String doi;
  final String url;
  final String found;
  final String why;
  final String taken;

  bool get scholarly =>
      doi.isNotEmpty ||
      _scholarlyKinds.any((scholarlyKind) => scholarlyKind.name == kind);

  MapSource copyWith({String? found, String? why, String? taken}) => MapSource(
    sourceId: sourceId,
    reference: reference,
    kind: kind,
    doi: doi,
    url: url,
    found: found ?? this.found,
    why: why ?? this.why,
    taken: taken ?? this.taken,
  );

  Map<String, Object?> toJson() => {
    'sourceId': sourceId,
    'reference': reference,
    'kind': kind,
    'doi': doi,
    'url': url,
    'found': found,
    'why': why,
    'taken': taken,
  };
}
