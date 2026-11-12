enum LayoutSize {
  compact,
  medium,
  expanded;

  static const mediumFrom = 600.0;
  static const expandedFrom = 1024.0;

  static LayoutSize of(double width) {
    if (width >= expandedFrom) {
      return LayoutSize.expanded;
    }
    if (width >= mediumFrom) {
      return LayoutSize.medium;
    }
    return LayoutSize.compact;
  }
}
