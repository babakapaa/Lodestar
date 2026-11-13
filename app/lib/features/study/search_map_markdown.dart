import '../../l10n/generated/app_localizations.dart';
import 'map_source.dart';
import 'search_map.dart';

class SearchMapMarkdown {
  const SearchMapMarkdown._();

  static String of(AppLocalizations l10n, String title, SearchMap map) {
    final out = StringBuffer('# ${l10n.mapTitle}: $title\n');
    void sources(String heading, List<MapSource> list, String why) {
      out.write('\n## $heading\n');
      if (list.isEmpty) {
        out.write('\n${l10n.mapNoneHere}\n');
      }
      for (final (index, source) in list.indexed) {
        out.write('\n${index + 1}. ${source.reference}\n');
        for (final (label, value) in [
          (l10n.mapFound, source.found),
          (why, source.why),
          (l10n.mapTaken, source.taken),
        ]) {
          if (value.trim().isNotEmpty) {
            out.write('   - $label: ${value.trim()}\n');
          }
        }
      }
    }

    sources(l10n.mapUsed, map.used, l10n.mapWhyUsed);
    sources(l10n.mapRejected, map.rejected, l10n.mapWhyRejected);
    out.write('\n## ${l10n.mapPath}\n\n${map.path.trim()}\n');
    if (map.queries.isNotEmpty) {
      out.write('\n');
      for (final query in map.queries) {
        out.write('- $query\n');
      }
    }
    return out.toString();
  }
}
