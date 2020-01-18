package task

import (
	"errors"
	"fmt"
	"strconv"

	"code.lampworld.xyz/go_lamp/websocket/config"
	msdb "code.lampworld.xyz/go_lamp/websocket/db"
)

//玩家

//Login 玩家登陆
func Login(a string, c string) (p PalyerInfo, err error) {
	Err := config.NewErr("Login")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("task.go: PlayerLogin err=", err)
		}
	}()
	for i := 0; i < len(tort.TeamInfo); i++ {
		if len(tort.TeamInfo[i].Palyers) > 0 {
			_, ok := tort.TeamInfo[i].Palyers[a]
			if ok && tort.TeamInfo[i].Palyers[a].Cipher == c {
				p = *tort.TeamInfo[i].Palyers[a]
				return
			}
		}
	}
	//数据库操作
	db := msdb.NewDb()
	if err = db.NewConn(); err != nil {
		return
	}
	defer db.Close()
	//查询数据库
	str := `select cipher,sdc,sdclv from Users where account='` + a + `'`
	res, err := db.Query(str)
	if err != nil {
		return
	}
	if len(res) != 1 {
		err = Err.Err6
		return
	}
	cipher := res[0]["cipher"]
	sdclv, err := strconv.ParseInt(res[0]["sdclv"], 10, 64)
	if err != nil {
		return
	}
	//验证密码
	if c == cipher {
		p = PalyerInfo{
			Account: a,
			Cipher:  c,
			Sdc:     0,
			Vip:     int(sdclv),
		}
		return
	}
	err = Err.Err5
	return
}

//JionTeam 加入队伍
func JionTeam(a string, c string, s float64, t int) (p PalyerInfo, err error) {
	Err := config.NewErr("Login")
	if tort.IsPutSDC == false {
		err = Err.Err3
		return
	}
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("task.go: PlayerLogin err=", err)
		}
	}()
	for i := 0; i < len(tort.TeamInfo); i++ {
		if len(tort.TeamInfo[i].Palyers) > 0 {
			_, ok := tort.TeamInfo[i].Palyers[a]
			if ok && tort.TeamInfo[i].Palyers[a].Cipher == c {
				p = *tort.TeamInfo[i].Palyers[a]
				if s == 0 {
					err = Err.Err4
					return
				} else {
					err = Err.Err1
					return
				}
			}
		}
	}
	//数据库操作
	db := msdb.NewDb()
	if err = db.NewConn(); err != nil {
		return
	}
	defer db.Close()
	//查询数据库
	str := `select cipher,sdc,sdclv from Users where account='` + a + `'`
	res, err := db.Query(str)
	if err != nil {
		return
	}
	if len(res) != 1 {
		err = Err.Err6
		return
	}
	cipher := res[0]["cipher"]
	sdclv, err := strconv.ParseFloat(res[0]["sdclv"], 64)
	if err != nil {
		return
	}
	if s > (sdclv+1)*10 {
		err = Err.Err4
		return
	}
	//验证密码
	if c == cipher {
		sdc, _ := strconv.ParseFloat(res[0]["sdc"], 64)
		if s > sdc {
			err = Err.Err9
			return
		}
		p, err = newPlayer(a, c, s, t, sdc)
		return
	}
	err = Err.Err5
	return
}

//玩家注册
func newPlayer(a string, c string, s float64, t int, oldsdc float64) (p PalyerInfo, err error) {
	Err := config.NewErr("Login")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("task.go: newPlayer err=", err)
		}
	}()
	p = PalyerInfo{
		Account: a,
		Cipher:  c,
		Sdc:     s,
		TeamID:  t,
	}
	//数据库操作
	db := msdb.NewDb()
	if err = db.NewConn(); err != nil {
		err = Err.Err7
		return
	}
	defer db.Close()
	//查询数据库
	str := fmt.Sprintf(`select * from Tortoise_History where users_id=(select id from users where account='%v') and Create_Date>'%v'`, a, getDay())
	res, err := db.Query(str)
	if err != nil {
		err = Err.Err7
		return
	}
	if len(res) < 1 {
		//加入玩家信息
		teamID := strconv.FormatInt(int64(t), 10)
		sdc := strconv.FormatFloat(s, 'f', 1, 64)
		str := `insert into Tortoise_History values(NEWID(),(select id from users where account='` + a + `'),(select ID  from Tortoise_Team where Team_ID=` + teamID + ` and Tortoise_ID =(select ID from Tortoise where Create_Date>'` + getDay() + `')),` + sdc + `,GETDATE(),0)`
		// fmt.Println(str)
		if _, err = db.Exec(str); err != nil {
			err = Err.Err7
			return
		}
		//神果操作记录
		str = fmt.Sprintf(`insert into Table_OperationHistory values((select id from users where account='%v'),'%v',100,GETDATE(),'赏金游戏扣除',%v,%v,0,0)`, a, a, oldsdc, oldsdc-s)
		// fmt.Println(str)
		if _, err = db.Exec(str); err != nil {
			err = Err.Err7
			return
		}
		//修改玩家神果数量
		str = `update users set sdc=sdc-` + sdc + ` where account='` + a + `'`
		if _, err = db.Exec(str); err != nil {
			err = Err.Err7
			return
		}
		tortChan <- &p
		tort.TeamInfo[t-1].TeamSDC += s
		tort.TeamInfo[t-1].Palyers[a] = <-tortChan
	}
	return
}

//生成玩家信息
func createPlayer(p PalyerInfo, t int) {
	tortChan <- &p
	tort.TeamInfo[t-1].TeamSDC += p.Sdc
	tort.TeamInfo[t-1].Palyers[p.Account] = <-tortChan
}

//Award 领取奖励
func (p *PalyerInfo) Award() error {
	Err := config.NewErr("Award")
	if !tort.IsGameEnd {
		return Err.Err3
	}
	if teamList[p.TeamID].Victory == false {
		return Err.Err10
	}
	if p.HasAward == true {
		return Err.Err1
	}
	//数据库操作
	db := msdb.NewDb()
	if err := db.NewConn(); err != nil {
		return Err.Err7
	}
	defer db.Close()

	//修改玩家记录
	str := fmt.Sprintf(`update Tortoise_History set has_Award=1 where Users_ID=(select ID from Users where Account='%v') and Create_Date>='%v'`, p.Account, getDay())

	if num, err := db.Exec(str); num < 1 || err != nil {
		return Err.Err7
	}

	//修改玩家神果数量
	str = fmt.Sprintf(`update Users set sdc=sdc+%.3f where account='%v'`, p.Sdc*teamList[p.TeamID].Scale, p.Account)

	if _, err := db.Exec(str); err != nil {
		return Err.Err7
	}

	//神果操作记录
	str = fmt.Sprintf(`insert into Table_OperationHistory values((select id from Users where Account='%v'),'%v',101,GETDATE(),'赏金任务获得',(select sdc from Users where Account='%v'),(select sdc from Users where Account='%v')+%v,0,0)`, p.Account, p.Account, p.Account, p.Account, p.Sdc)
	// fmt.Println(str)
	if _, err := db.Exec(str); err != nil {
		return Err.Err7
	}

	(*p).HasAward = true
	errStr := fmt.Sprintf(`{"LoginResult":{"code":0,"msg":"成功","sdc":%.2f}}`, p.Sdc*teamList[p.TeamID].Scale)
	return errors.New(errStr)
}
