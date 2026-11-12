abstract final class PointsText {
  static String of(double value) =>
      value == value.roundToDouble() ? '${value.toInt()}' : '$value';
}
