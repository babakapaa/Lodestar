import 'package:lodestar/core/files/file_opener.dart';
import 'package:lodestar/core/files/opened_file.dart';

class FakeFileOpener implements FileOpener {
  List<OpenedFile> next = const [];

  @override
  Future<List<OpenedFile>> open(List<String> extensions) async => next;
}
