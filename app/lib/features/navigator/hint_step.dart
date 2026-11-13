enum HintStep {
  question,
  direction,
  source,
  fragment,
  example;

  static HintStep? parse(Object? value) {
    for (final step in values) {
      if (step.name == value) {
        return step;
      }
    }
    return null;
  }
}
