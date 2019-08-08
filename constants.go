package nxrm

const (
	ConfFileName = "nexus3-repository-cli.json"
	connDetailsSuccessInfo = "Connection details were stored successfully in the file ./%s\n"
	connDetailsEmptyInfo   = "Server connection details are not set...First Run %q to set the connection details\n"

	// API Extensions
	apiBase        = "service/rest"
	scriptAPI      = "v1/script"
	repositoryPath = "v1/repositories"

	successStatus   = "200 OK"
	notFoundStatus  = "404 Not Found"
	noContentStatus = "204 No Content"
	foundStatus     = "302 Found"

	// Script Path
	scriptBasePath = "./scripts/groovy"

	// Error Strings
	jsonMarshalError   = "JSON Marshal Error"
	jsonUnmarshalError = "JSON Unmarshal Error"
	setVerboseInfo   = "There was an error calling the function. Set verbose flag for more information"

	nameRequiredInfo = "name is a required parameter"

	//script
	scriptAddedInfo       = "The script %q is added to nexus\n"
	scriptUpdatedInfo     = "The script %q is updated in nexus\n"
	scriptDeletedInfo     = "The script %q is deleted from nexus\n"
	scriptRunSuccessInfo  = "The script %q was executed successfully\n"
	scriptRunNotFoundInfo = "The script %q was not found in nexus. Make sure you add the script to nexus before executing the script\n"
	scriptExistsInfo      = "The script %q already exists in nexus\n"
	scriptNotfoundInfo    = "The script %q was not found in nexus\n"

	//scripts
	getRepoScript            = "get-repo"
	createHostedRepoScript   = "create-hosted-repo"
	createProxyRepoScript    = "create-proxy-repo"
	createGroupRepoScript    = "create-group-repo"
	updateGroupMembersScript = "update-group-members"
	deleteRepoScript         = "delete-repo"
	getPrivilegesScript      = "get-privileges"
	createPrivilegeScript    = "create-privilege"
	updatePrivilegeScript    = "update-privilege"
	deletePrivilegeScript    = "delete-privilege"
	getRoleScript            = "get-roles"
	createRoleScript         = "create-role"
	deleteRoleScript         = "delete-role"

	//repo
	RepoFormatNotValidInfo        = "%q is not a valid repository format. Available repository formats are : %v\n"
	repoFormatRequiredInfo        = "format is a required parameter"
	hostedRepoRequiredInfo        = "name and format are required parameters to create a hosted repository"
	proxyRepoRequiredInfo         = "name, format and remote-url are required parameters to create a proxy repository"
	groupRequiredInfo             = "name, format and members are required parameters"
	dockerPortsInfo               = "You need to specify either a http port or a https port or both for creating a docker repository"
	repositoryNotFoundInfo        = "Repository %q was not found in nexus"
	repoCreatedInfo               = "Repository %q was created in nexus\n"
	repoUpdatedStatus             = "Repository %q was updated in nexus\n"
	repoDeletedInfo               = "Repository %q was deleted from nexus\n"
	repoCreateErrorInfo           = "Error creating repository : %s\n"
	repoUpdateErrorInfo           = "Error updating repository : %s\n"
	repoDeleteErrorInfo           = "Error deleting repository : %s\n"
	repoExistsInfo                = "Repository %q already exists in nexus\n"
	cannotBeSameRepoInfo          = "Member %q == group %q, cannot add a group repository as a member in the same group\n"
	proxyCredsNotValidInfo        = "You need to provide both proxy-user and proxy-pass to set credentials to a proxy repository"
	remoteURLNotValidInfo         = "%q is an invalid url. URL must begin with either http:// or https://"
	notAGroupRepoInfo             = "%q is not a group repository\n"
	groupMemberInvalidFormatInfo  = "Repository %q is not a %q format repository, hence it cannot be added to the group repository\n"
	groupMemberAlreadyExistsInfo  = "Member %q already exists in the group %q, hence not adding the member again\n"
	groupMemberNotFoundInfo       = "Repository %q was not found in Nexus, hence it cannot be added to the group repository\n"
	groupMemberRemoveNotFoundInfo = "Member %q was not found in the group %q, hence cannot remove the member from the group\n"
	groupMemberRequiredInfo       = "At least one valid group member should be provided to add to a group repository"
	groupMemberAddSuccessInfo     = "Member %q is added to the group %q\n"
	groupMemberRemoveSuccessInfo  = "Member %q is removed from the group %q\n"

	//selector
	contentSelectorType = "csel"
	defaultContentSelectorDescription = "Custom content-selector created from the CLI"
	createSelectorRequiredInfo        = "name and expression are required parameters"
	createSelectorSuccessInfo         = "Content selector %q was created\n"
	updateSelectorSuccessInfo         = "Content selector %q was updated\n"
	deleteSelectorSuccessInfo         = "Content selector %q was deleted\n"
	selectorAlreadyExistsInfo         = "Content selector %q already exists in nexus\n"
	selectorNotFoundInfo              = "Content selector %q was not found in nexus\n"

	//privilege
	defaultPrivilegeDescription = "Custom privilege created from the CLI"
	privilegeNotFoundInfo       = "Privilege %q was not found in nexus\n"
	privilegeExistsInfo         = "Privilege %q already exists\n"
	createPrivilegeRequiredInfo = "name, selector-name and repo-name are required parameters"
	createPrivilegeSuccessInfo  = "Privilege %q is created"
	updatePrivilegeSuccessInfo  = "Privilege %q is updated"
	deletePrivilegeSuccessInfo  = "Privilege %q is deleted"

	//role
	UpdateActionRequiredInfo = "Update action is a required parameter. Available values = %+q\n"
	UpdateActionInvalidInfo  = "%s is not a valid update action. Available actions: %+q\n"

	defaultRoleDescription        = "Custom role created from the CLI"
	defaultRoleSource             = "Nexus"
	roleIDRequiredInfo            = "id is a required parameter"
	roleNotFoundInfo              = "Role %q was not found in nexus\n"
	roleExistsInfo                = "Role %q already exists\n"
	createRoleRequiredInfo        = "id, description and source are required parameters"
	createRoleSuccessInfo         = "Role %q is created with role members %v and privileges %+q\n"
	updateRoleSuccessInfo         = "Role %q is updated\n"
	deleteRoleSuccessInfo         = "Role %q is deleted\n"
	roleItemsRequiredInfo         = "%s : You need to provide at least one valid role member or role privilege during role creation\n"
	noRoleMemberProvidedInfo      = "No role members are provided to add/remove to/from the role"
	noValidRoleMemberInfo         = "No valid role members are provided to add to the role"
	cannotBeSameRoleInfo          = "Role member %q == role id %s, cannot add a role as a member in the same role"
	roleMemberNotFoundInfo        = "Role %q was not found in nexus, hence it cannot be added to the role"
	rolePrivilegeNotFoundInfo     = "Privilege %q was not found in nexus, hence it cannot be added to the role"
	noValidRolePrivilegeInfo      = "No valid privileges are provided to add to the role"
	noRolePrivilegesIProvidedInfo = "No privileges are provided to add to the role"
)
