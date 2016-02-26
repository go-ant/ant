package data

type rolePermissions struct {
	Name        string
	Slug        string
	Description string
	Permissions map[string][]string
}

type post struct {
	Title    string
	Slug     string
	Markdown string
}
