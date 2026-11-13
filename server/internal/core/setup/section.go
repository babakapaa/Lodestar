package setup

type Section struct {
	Key      string
	Required bool
}

const (
	SectionDeploy       = "deploy"
	SectionUniversity   = "university"
	SectionPeople       = "people"
	SectionAdmins       = "admins"
	SectionIntegrations = "integrations"
	SectionTeaching     = "teaching"
	SectionNavigator    = "navigator"
	SectionLeague       = "league"
	SectionSandbox      = "sandbox"
	SectionCampus       = "campus"
	SectionPrivacy      = "privacy"
	SectionLegal        = "legal"
	SectionLaunch       = "launch"
)

var Sections = []Section{
	{Key: SectionDeploy, Required: true},
	{Key: SectionUniversity, Required: true},
	{Key: SectionPeople, Required: true},
	{Key: SectionAdmins, Required: true},
	{Key: SectionIntegrations},
	{Key: SectionTeaching},
	{Key: SectionNavigator},
	{Key: SectionLeague},
	{Key: SectionSandbox},
	{Key: SectionCampus},
	{Key: SectionPrivacy},
	{Key: SectionLegal, Required: true},
	{Key: SectionLaunch},
}
