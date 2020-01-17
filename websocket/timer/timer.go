package timer

import (
	"fmt"
	"go_code/websocket/config"
	"go_code/websocket/task"
	"time"
)

var (
	d int //日
	h int //时
	m int //分
)

//Timer 定时器
func Timer() {
	fmt.Println(getTime(), "定时任务启动...")
	tort := task.GetTort()
	cTime := config.NewTask()
	for {
		d = time.Now().Day()
		h = time.Now().Hour()
		m = time.Now().Minute()
		//每日任务
		if d > cTime.T1 {
			fmt.Println(getTime(), "开始新的一天...")
			if tort.IsGameEnd {
				fmt.Println(getTime(), "清除比赛信息...")
				*tort = task.Tortoise{
					IsPutSDC:    true,
					IsGameStart: false,
					IsGameEnd:   false,
					GameResult:  `{"0":0,"1":0,"2":0}`,
					TeamInfo:    nil,
				}
			}
			if !tort.IsPutSDC {
				fmt.Println(getTime(), "开启押注...")
				tort.IsPutSDC = true
			}
			fmt.Println(getTime(), "创建乌龟表...")
			if err := task.CreateTable(); err != nil {
				fmt.Println(err)
			}
			fmt.Println(getTime(), "创建队伍表...")
			if err := task.CreateTeam(); err != nil {
				fmt.Println(err)
			}
			fmt.Println(getTime(), "备份数据库...")
			if err := task.BackupDb(); err != nil {
				fmt.Println("备份数据库失败")
			}
			cTime.T1 = d
		}

		//定时任务
		if h == cTime.T2-1 && m == 58 {
			if tort.IsPutSDC {
				fmt.Println(getTime(), "关闭押注...")
				tort.IsPutSDC = false
			}
		}
		if h == cTime.T2 {
			if m < 1 && !tort.IsGameStart {
				fmt.Println(getTime(), "开始游戏...")
				tort.IsGameStart = true
				tort.GameStart()
			}
			if m > 2 && tort.IsGameStart {
				fmt.Println(getTime(), "游戏结束...")
				tort.IsGameStart = false
			}
			if m >= 5 && !tort.IsGameEnd {
				fmt.Println(getTime(), "开启奖励领取...")
				tort.IsGameEnd = true
			}
		}
		time.Sleep(time.Second)
	}
}

//获取时间 2006-01-02 15:04:05
func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
