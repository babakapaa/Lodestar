package audit

const ViewAudit = "core.audit.view"

func Permissions() []string {
	return []string{ViewAudit}
}
