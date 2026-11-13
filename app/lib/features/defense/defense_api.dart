import '../../core/api/api_client.dart';
import '../study/grade_input.dart';
import 'defense.dart';
import 'defense_question.dart';
import 'defense_row.dart';
import 'defense_slot.dart';
import 'defense_status.dart';
import 'my_defense.dart';
import 'new_slots.dart';

class DefenseApi {
  const DefenseApi(this._api);

  final ApiClient _api;

  Future<List<DefenseSlot>> slots(String assignmentId) async => [
    for (final item
        in (await _api.get('${_work(assignmentId)}/defense-slots'))!
            as List<Object?>)
      DefenseSlot.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> addSlots(String assignmentId, NewSlots slots) =>
      _api.post('${_work(assignmentId)}/defense-slots', slots.toJson());

  Future<void> deleteSlot(String id) =>
      _api.delete('/v1/study/defense-slots/${Uri.encodeComponent(id)}');

  Future<MyDefense> mine(String assignmentId) async => MyDefense.fromJson(
    (await _api.get('${_work(assignmentId)}/defense'))! as Map<String, Object?>,
  );

  Future<MyDefense> book(String assignmentId, String slotId) async =>
      MyDefense.fromJson(
        (await _api.put('${_work(assignmentId)}/defense/booking', {
              'slotId': slotId,
            }))!
            as Map<String, Object?>,
      );

  Future<MyDefense> cancel(String assignmentId) async => MyDefense.fromJson(
    (await _api.delete('${_work(assignmentId)}/defense/booking'))!
        as Map<String, Object?>,
  );

  Future<List<DefenseRow>> list(String assignmentId) async => [
    for (final item
        in (await _api.get('${_work(assignmentId)}/defenses'))!
            as List<Object?>)
      DefenseRow.fromJson(item! as Map<String, Object?>),
  ];

  Future<void> sample(String assignmentId, int count) =>
      _api.post('${_work(assignmentId)}/defenses/sample', {'count': count});

  Future<Defense> get(String assignmentId, String studentId) async =>
      Defense.fromJson(
        (await _api.get(_student(assignmentId, studentId)))!
            as Map<String, Object?>,
      );

  Future<Defense> save(
    String assignmentId,
    String studentId,
    List<DefenseQuestion> questions,
    String summary,
  ) async => Defense.fromJson(
    (await _api.put(
          _student(assignmentId, studentId),
          _body(questions, summary),
        ))!
        as Map<String, Object?>,
  );

  Future<void> select(String assignmentId, String studentId, bool selected) =>
      _api.put('${_student(assignmentId, studentId)}/selection', {
        'selected': selected,
      });

  Future<void> finish(
    String assignmentId,
    String studentId,
    DefenseStatus status,
    GradeInput? grade,
  ) => _api.post('${_student(assignmentId, studentId)}/finish', {
    'status': status.name,
    'grade': grade?.toJson(),
  });

  Future<List<DefenseQuestion>> draftQuestions(
    String assignmentId,
    String studentId,
  ) async {
    final body =
        (await _api.post(
              '${_student(assignmentId, studentId)}/questions-draft',
            ))!
            as Map<String, Object?>;
    return [
      for (final item in body['questions']! as List<Object?>)
        DefenseQuestion.fromJson(item! as Map<String, Object?>),
    ];
  }

  Future<String> draftProtocol(
    String assignmentId,
    String studentId,
    List<DefenseQuestion> questions,
    String summary,
  ) async =>
      ((await _api.post(
                '${_student(assignmentId, studentId)}/protocol-draft',
                _body(questions, summary),
              ))!
              as Map<String, Object?>)['summary']
          as String? ??
      '';

  static Map<String, Object?> _body(
    List<DefenseQuestion> questions,
    String summary,
  ) => {
    'questions': [for (final q in questions) q.toJson()],
    'summary': summary,
  };

  static String _work(String id) =>
      '/v1/study/assignments/${Uri.encodeComponent(id)}';

  static String _student(String assignmentId, String studentId) =>
      '${_work(assignmentId)}/defenses/${Uri.encodeComponent(studentId)}';
}
