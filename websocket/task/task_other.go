package task

import (
	"errors"
	msdb "go_code/websocket/db"
	"strconv"
	"time"
)

//其他

//同步数据库信息
func checkMsdb() error {

	var (
		str string
		res []map[string]string
		err error
	)

	//数据库操作
	db := msdb.NewDb()
	db.NewConn()
	defer db.Close()

	//查询乌龟历史记录 获取玩家 ID 列表
	str = `select CAST(Users_ID as varchar(100)),sdc,has_award from Tortoise_History where create_date>='` + getDay() + `'`
	if res, _ = db.Query(str); len(res) < 1 {
		return errors.New("task_other.go: checkMsdb >获取玩家 ID 列表失败")
	}

	//遍历玩家 ID 列表 获取用户 账号、密码
	for _, v := range res {
		str = `select account,cipher,sdclv from users where id='` + v[""] + `'`
		if res, _ = db.Query(str); len(res) < 1 {
			return errors.New("task_other.go: checkMsdb >获取用户 账号、密码 失败")
		}

		//获取个人信息
		p := res[0]

		//获取玩家队伍 ID
		str = `select Team_ID from Tortoise_Team where id=(select top 1 Team_ID from Tortoise_History where users_id='` + v[""] + `')`
		if res, _ = db.Query(str); len(res) != 1 {
			return errors.New("task_other.go: checkMsdb >获取玩家队伍 ID 失败")
		}
		var n int
		if n, err = strconv.Atoi(res[0]["Team_ID"]); err != nil {
			return errors.New("task_other.go: checkMsdb >字符串转换 int 失败")
		}

		//数据转换
		var sdc float64
		sdc, err = strconv.ParseFloat(v["sdc"], 64)
		var vip int64
		vip, err = strconv.ParseInt(p["sdclv"], 10, 64)
		if err != nil {
			return errors.New("task_other.go: checkMsdb >字符串转换 失败")
		}

		//判断是否领取奖励
		var award bool
		if v["has_award"] == "0" {
			award = false
		} else {
			award = true
		}

		//生成玩家对象
		player := PalyerInfo{
			Account:  p["account"],
			Cipher:   p["cipher"],
			Sdc:      sdc,
			Vip:      int(vip),
			HasAward: award,
			TeamID:   n,
		}

		//创建玩家
		createPlayer(player, n)
	}

	return nil
}

//获取日期
func getDay() string {
	return time.Now().Format("2006-01-02")
}

//BackupDb 备份数据库
func BackupDb() error {

	//数据库操作
	db := msdb.NewDb()
	db.NewConn()
	defer db.Close()

	//数据库备份
	str := `Backup Database ShenDengWord To disk='D:/Backup/` + getDay() + `.bak'`
	if _, err := db.Query(str); err != nil {
		return errors.New("task_other.go: BackupDb >备份数据库失败")
	}
	return nil
}

//mapSort 排序
func mapSort(m map[int]float64) int {
	var (
		res float64
		key int
	)
	for k, v := range m {
		if res < v {
			key = k
			res = v
		}
	}
	return key
}
