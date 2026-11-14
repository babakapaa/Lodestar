import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'tester_person.dart';

class TesterNext extends Notifier<TesterPerson?> {
  @override
  TesterPerson? build() => null;

  void set(TesterPerson? person) => state = person;

  TesterPerson? take() {
    final person = state;
    state = null;
    return person;
  }
}
