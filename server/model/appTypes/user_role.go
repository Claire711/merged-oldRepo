package appTypes

type RoleID int

const (
	//iota：自动枚举器
	Guest RoleID = iota //游客 iota=0 → Guest=0
	User                //用户 iota=1 → User=1
	Admin               //管理员 iota=2 → Admin=2
)
