package task

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"code.lampworld.xyz/go_lamp/websocket/config"
	msdb "code.lampworld.xyz/go_lamp/websocket/db"
)

//队伍

//CreateTeam 创建队伍
func CreateTeam() error {
	t1, err := newTorTeam(1)
	t2, err := newTorTeam(2)
	t3, err := newTorTeam(3)
	if err != nil {
		return err
	}
	teamList[1] = t1
	teamList[2] = t2
	teamList[3] = t3
	return nil
}

//newTorTeam 生成队伍信息
func newTorTeam(id int) (*TortoiseTeam, error) {
	team := TortoiseTeam{
		TeamID:  id,
		TeamSDC: 0,
		Victory: false,
		Scale:   0,
		Palyers: make(map[string]*PalyerInfo, 0),
	}
	//数据库操作
	db := msdb.NewDb()
	if err := db.NewConn(); err != nil {
		return nil, err
	}
	defer db.Close()
	//查询数据库
	str := `select Team_ID,Team_SDC,Team_Num from Tortoise_Team where Tortoise_ID=(select id from Tortoise where Create_Date>'` + getDay() + `') and Team_ID=` + strconv.FormatInt(int64(id), 10)
	res, err := db.Query(str)
	if len(res) < 1 {
		//修改数据库
		str := `insert into Tortoise_Team values((select id from Tortoise where Create_Date>'` + getDay() + `'),` + strconv.FormatInt(int64(id), 10) + `,0,0)`
		if _, err := db.Exec(str); err != nil {
			return nil, err
		}
		tort.TeamInfo = append(tort.TeamInfo, &team)
		return &team, nil
	}
	teamID, _ := strconv.ParseInt(res[0]["Team_ID"], 10, 0)
	teamSDC, _ := strconv.ParseFloat(res[0]["Team_SDC"], 64)
	team.TeamID = int(teamID)
	team.TeamSDC = teamSDC
	tort.TeamInfo = append(tort.TeamInfo, &team)
	return &team, err
}

//比赛

//GameStart 开始比赛
func (t *Tortoise) GameStart() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("task.go: GameStart err=", err)
		}
	}()
	//获取队伍信息
	//获胜概率:((队伍数量-1)-((所在队伍神果数/神果总数)/(队伍数量-1)+(所在队伍人数/总人数)/(队伍数量-1)))/(队伍数量-1)
	num := len(t.TeamInfo) //队伍数量
	var (
		sdc    float64 //神果总数
		player float64 //总人数
	)
	sdc = 0
	player = 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		sdc += t.TeamInfo[i].TeamSDC
		player += float64(len(t.TeamInfo[i].Palyers))
	}
	p := make(map[int]float64, 3)
	if sdc == 0 || player == 0 {
		for i := 0; i < num; i++ {
			p[i] = rand.Float64()
		}
	} else {
		for i := 0; i < num; i++ {
			r := rand.Float64() / float64(num)
			pSdc := (1 - t.TeamInfo[i].TeamSDC/sdc) / (float64(num) - 1)
			pPlayer := (1 - float64(len(t.TeamInfo[i].Palyers))/float64(player)) / (float64(num) - 1)
			if t.TeamInfo[i].TeamSDC > 0 {
				p[i] = (pSdc+pPlayer)/2 + r
			} else {
				p[i] = 0
			}
		}
	}
	vid := mapSort(p)
	t.TeamInfo[vid].Victory = true
	t.TeamInfo[vid].Scale = 1 + ((1 - config.NewTask().Charge) * (sdc - t.TeamInfo[vid].TeamSDC) / sdc)
	scale := 1 / p[vid]
	//生成比赛情况
	m1 := rand.Intn(3)
	m2 := rand.Intn(3)
	m3 := rand.Intn(3)
	go func() {
		for i := 1; i <= 1500; i++ {
			second := float64(i)
			z := make(map[int]float64)
			z[0] = createGameResult(p[0], second, 1500, m1+1, scale)
			z[1] = createGameResult(p[1], second, 1500, m2+1, scale)
			z[2] = createGameResult(p[2], second, 1500, m3+1, scale)
			str := fmt.Sprintf(`{"0":%.6f,"1":%.6f,"2":%.6f}`, z[0], z[1], z[2])
			tort.GameResult = str
			time.Sleep(time.Millisecond * 40)
		}
	}()
}
