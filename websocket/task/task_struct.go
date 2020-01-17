package task

//PalyerInfo 玩家信息
type PalyerInfo struct {
	//玩家账号
	Account string
	//玩家密码
	Cipher string
	//投入的神果数
	Sdc float64
	//vip 等级
	Vip int
	//领取奖励
	HasAward bool
	//队伍编号
	TeamID int
}

//Tortoise 游戏信息
type Tortoise struct {
	IsPutSDC    bool            //是否允许押注
	IsGameStart bool            //游戏是否开始
	IsGameEnd   bool            //是否允许结算
	GameResult  string          //比赛结果
	TeamInfo    []*TortoiseTeam //团队信息
}

//TortoiseTeam 队伍信息
type TortoiseTeam struct {
	TeamID  int                    //团队编号
	TeamSDC float64                //神果数
	Victory bool                   //是否获胜
	Scale   float64                //sdc 比例
	Palyers map[string]*PalyerInfo //用户列表
}

type moveLine struct {
	f1 func(t float64, S float64, T float64) float64
	f2 func(t float64, S float64, T float64) float64
	f3 func(t float64, S float64, T float64) float64
}
