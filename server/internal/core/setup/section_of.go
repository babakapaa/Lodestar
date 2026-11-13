package setup

import "strings"

var exactSections = map[string]string{
	"core.default_language":      SectionUniversity,
	"core.default_theme":         SectionUniversity,
	"core.time_zone":             SectionUniversity,
	"core.lti_assignment_source": SectionIntegrations,
	"core.audit_retention_days":  SectionPrivacy,
	"core.self_deletion":         SectionPrivacy,
	"core.deletion_grace_days":   SectionPrivacy,
	"core.files.max_file_mb":     SectionPrivacy,
	"chat.retention_days":        SectionPrivacy,
	"study.defense_recording":    SectionLegal,
	"chat.call_recording":        SectionLegal,
	"chat.recording_consent":     SectionLegal,
	"navigator.support_contacts": SectionLegal,
	"sandbox.watch_processes":    SectionLegal,
}

var prefixSections = map[string]string{
	"study":      SectionTeaching,
	"team":       SectionTeaching,
	"navigator":  SectionNavigator,
	"league":     SectionLeague,
	"sandbox":    SectionSandbox,
	"attendance": SectionCampus,
	"chat":       SectionCampus,
}

func SectionOf(key string) string {
	if s, ok := exactSections[key]; ok {
		return s
	}
	prefix, _, _ := strings.Cut(key, ".")
	if s, ok := prefixSections[prefix]; ok {
		return s
	}
	return SectionUniversity
}
