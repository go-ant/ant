package models

type Menu struct {
	Name  string
	URL   string
	Label string
	Icon  string
	List  []Menu
}
