import 'map_gap.dart';
import 'map_rules.dart';
import 'map_source.dart';

class SearchMap {
  const SearchMap({
    this.used = const [],
    this.rejected = const [],
    this.path = '',
    this.queries = const [],
  });

  factory SearchMap.fromJson(Map<String, Object?> json) => SearchMap(
    used: MapSource.list(json['used']),
    rejected: MapSource.list(json['rejected']),
    path: json['path'] as String? ?? '',
    queries: [
      for (final item in json['queries'] as List<Object?>? ?? const [])
        item! as String,
    ],
  );

  static SearchMap? maybe(Object? json) =>
      json is Map<String, Object?> ? SearchMap.fromJson(json) : null;

  final List<MapSource> used;
  final List<MapSource> rejected;
  final String path;
  final List<String> queries;

  bool get isEmpty =>
      used.isEmpty &&
      rejected.isEmpty &&
      path.trim().isEmpty &&
      queries.isEmpty;

  int get scholarlyShare => used.isEmpty
      ? 0
      : (used.where((s) => s.scholarly).length * 100 / used.length).floor();

  List<MapGap> gaps(MapRules rules) => [
    if (used.isEmpty || used.length < rules.minSources) MapGap.sources,
    if (used.isNotEmpty &&
        used.where((s) => s.scholarly).length * 100 <
            rules.minScholarly * used.length)
      MapGap.scholarly,
    if ([...used, ...rejected].any((s) => s.why.trim().isEmpty)) MapGap.why,
    if (path.trim().isEmpty) MapGap.path,
  ];

  bool contains(String sourceId) =>
      [...used, ...rejected].any((s) => s.sourceId == sourceId);

  SearchMap copyWith({
    List<MapSource>? used,
    List<MapSource>? rejected,
    String? path,
    List<String>? queries,
  }) => SearchMap(
    used: used ?? this.used,
    rejected: rejected ?? this.rejected,
    path: path ?? this.path,
    queries: queries ?? this.queries,
  );

  Map<String, Object?> toJson() => {
    'used': [for (final s in used) s.toJson()],
    'rejected': [for (final s in rejected) s.toJson()],
    'path': path,
    'queries': queries,
  };
}
