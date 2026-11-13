package files

func ObjectKey(tenantID, fileID string) string {
	return tenantID + "/" + fileID
}
