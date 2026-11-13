enum MapGap {
  sources,
  scholarly,
  why,
  path;

  static List<MapGap> list(Object? json) => [
    for (final item in json as List<Object?>? ?? const [])
      for (final gap in values)
        if (gap.name == item) gap,
  ];
}
