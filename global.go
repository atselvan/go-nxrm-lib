package nxrm

// AuthUser represents the credential for Authentication
type AuthUserStruct struct {
	Username string
	Password string
}

type RequestBody struct {
	Json []byte
	Text string
}

var (
	NexusURL            string
	AuthUser            AuthUserStruct
	Verbose             bool
	Debug               bool
	SkipTLSVerification bool

	InitialRepoList  = []string{"maven-public", "maven-central", "maven-snapshots", "maven-releases", "nuget-group", "nuget-hosted", "nuget.org-proxy"}
	NexusScripts     = []string{"get-repo", "create-hosted-repo", "create-proxy-repo", "create-group-repo", "update-group-members", "delete-repo", "get-content-selectors", "create-content-selector", "update-content-selector", "delete-content-selector", "get-privileges", "create-privilege", "update-privilege", "delete-privilege", "get-roles", "create-role", "delete-role"}
	RepoFormats      = []string{"maven", "npm", "nuget", "bower", "pypi", "raw", "rubygems", "yum", "docker"}
	RepoType         = []string{"hosted", "proxy", "group"}
	PrivilegeActions = []string{"read", "write"}
	UpdateActions    = []string{"add", "remove"}
)
