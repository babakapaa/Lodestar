enum NavigatorIntent {
  understand('understand'),
  find('find'),
  check('check'),
  doForMe('do_for_me'),
  organize('organize'),
  state('state'),
  closed('closed');

  const NavigatorIntent(this.wire);

  final String wire;

  static const choosable = [understand, find, check, doForMe, organize, state];

  static NavigatorIntent? parse(Object? value) {
    for (final intent in values) {
      if (intent.wire == value) {
        return intent;
      }
    }
    return null;
  }
}
