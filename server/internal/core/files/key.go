package files

func objectKey(tenantID, fileID string) string {
	return tenantID + "/" + fileID
}
