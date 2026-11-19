import 'package:flutter_test/flutter_test.dart';
import 'package:lodestar/features/study/assignment.dart';
import 'package:lodestar/features/study/grade_parts.dart';
import 'package:lodestar/features/study/grade_weights.dart';

void main() {
  test('the grade is added up from weighted parts', () {
    const parts = GradeParts(defense: 80, map: 100, work: 50);
    expect(GradeWeights.suggested.points(parts, 10), 8);
    expect(
      GradeWeights.suggested.points(const GradeParts(defense: 80), 10),
      isNull,
    );
    expect(parts.toJson(), {'defense': 80.0, 'map': 100.0, 'work': 50.0});
  });

  test('weights must add up to 100', () {
    expect(GradeWeights.suggested.valid, isTrue);
    expect(const GradeWeights(defense: 50, map: 30, work: 10).valid, isFalse);
  });

  test('an assignment reads its weights', () {
    final assignment = Assignment.fromJson({
      'id': 'a1',
      'courseId': 'c1',
      'title': 'Эссе',
      'grading': 'points',
      'maxPoints': 10,
      'gradeWeights': {'defense': 50, 'map': 30, 'work': 20},
    });
    expect(assignment.gradeWeights?.map, 30);
  });
}
