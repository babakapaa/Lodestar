import 'opened_file.dart';

abstract interface class FileOpener {
  Future<List<OpenedFile>> open(List<String> extensions);
}
