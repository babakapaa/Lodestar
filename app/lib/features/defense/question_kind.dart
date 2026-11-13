enum QuestionKind {
  understanding,
  analysis,
  application,
  path;

  static QuestionKind parse(Object? value) => switch (value) {
    'analysis' => analysis,
    'application' => application,
    'path' => path,
    _ => understanding,
  };
}
