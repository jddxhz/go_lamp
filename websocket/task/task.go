package task

import (
	"fmt"
	msdb "go_code/websocket/db"
	// "strcov"
	"time"
)

var (
	tort = Tortoise{
		IsPutSDC:    true,
		IsGameStart: false,
		IsGameEnd:   false,
		GameResult:  `{"0":0,"1":0,"2":0}`,
		TeamInfo:    nil,
	}
	tortChan = make(chan *PalyerInfo, 10)
	teamList = make(map[int]*TortoiseTeam, 3)
)

var moveL = moveLine{
	f1: func(t float64, S float64, T float64) float64 {
		v := S / T
		s := v * t
		return s
	},
	f2: func(t float64, S float64, T float64) float64 {
		a := 2 * S / (T * T)
		s := 0.5 * a * t * t
		return s
	},
	f3: func(t float64, S float64, T float64) float64 {
		a := 4 * S / (T * T)
		if t < T*0.5 {
			s := 0.5 * a * t * t
			return s
		}
		s := S - 0.5*a*(T-t)*(T-t)
		return s
	},
}

//GetTort 获取比赛信息
func GetTort() *Tortoise {
	return &tort
}

//createGameResult 生成比赛结果
func createGameResult(p float64, t float64, T float64, m int, s float64) float64 {
	switch m {
	case 1:
		p = moveL.f1(t, p, T) * s
	case 2:
		p = moveL.f2(t, p, T) * s
	case 3:
		p = moveL.f3(t, p, T) * s
	}
	return p
}

//TortInfo 发给客户端的队伍信息
func (t Tortoise) TortInfo() string {
	str := make(map[int]string, 3)
	for i := 0; i < len(t.TeamInfo); i++ {
		a := len(t.TeamInfo[i].Palyers)
		b := t.TeamInfo[i].TeamSDC
		str[i] = fmt.Sprintf(`"%v":{"0":%v,"1":%v}`, i, a, b)
	}
	res := `{` + str[0] + "," + str[1] + "," + str[2] + `}`
	res = fmt.Sprintf(`{"TortInfo":{"isPutSDC":%v,"isGameStart":%v,"isGameEnd":%v,"teamInfo":%v,"gameResult":%v,"time":"%v"}}`, t.IsPutSDC, t.IsGameStart, t.IsGameEnd, res, tort.GameResult, time.Now().Format("2006-01-02 15:04:05"))
	return res
}

//CreateTable 创建乌龟表
func CreateTable() error {
	//数据库操作
	db := msdb.NewDb()
	if err := db.NewConn(); err != nil {
		return err
	}
	defer db.Close()

	date := time.Now().Format("2006-01-02")
	//查询数据库
	str := `select id from Tortoise where Create_Date>='` + date + `'`
	res, err := db.Query(str)
	if err != nil {
		return err
	}
	if len(res) < 1 {
		//修改数据库
		str = `insert into Tortoise values(0.1,'` + date + ` 20:00:00',GETDATE())`
		_, err := db.Exec(str)
		return err
	}
	return err
}

// 初始化
func init() {
	//创建乌龟表
	if err := CreateTable(); err != nil {
		fmt.Println(err)
	}
	if err := CreateTeam(); err != nil {
		fmt.Println(err)
	}
	if err := checkMsdb(); err != nil {
		fmt.Println(err)
	}
}
