package config

import (
	"errors"
	"time"
)

//TaskTime 任务时间
type TaskTime struct {
	T1     int     //日期
	T2     int     //小时
	Charge float64 //手续费
}

type err struct {
	Err0  error
	Err1  error
	Err2  error
	Err3  error
	Err4  error
	Err5  error
	Err6  error
	Err7  error
	Err8  error
	Err9  error
	Err10 error
}

//NewErr 错误信息列表
func NewErr(str string) err {
	return err{
		Err0:  errors.New(`{"` + str + `Result":{"code":0,"msg":"成功"}}`),
		Err1:  errors.New(`{"` + str + `Result":{"code":-1,"msg":"只能操作一次"}}`),
		Err2:  errors.New(`{"` + str + `Result":{"code":-2,"msg":"还没有投入神果"}}`),
		Err3:  errors.New(`{"` + str + `Result":{"code":-3,"msg":"错误的时间"}}`),
		Err4:  errors.New(`{"` + str + `Result":{"code":-4,"msg":"意料之外的神果数量"}}`),
		Err5:  errors.New(`{"` + str + `Result":{"code":-5,"msg":"密码错误"}}`),
		Err6:  errors.New(`{"` + str + `Result":{"code":-6,"msg":"用户不存在"}}`),
		Err7:  errors.New(`{"` + str + `Result":{"code":-7,"msg":"查询数据失败"}}`),
		Err8:  errors.New(`{"` + str + `Result":{"code":-8,"msg":"不可思议的错误"}}`),
		Err9:  errors.New(`{"` + str + `Result":{"code":-9,"msg":"神果数量不足"}}`),
		Err10: errors.New(`{"` + str + `Result":{"code":-10,"msg":"不符合领取条件"}}`),
	}
}

var (
	taskTime *TaskTime
)

//NewTask 任务时间
func NewTask() *TaskTime {
	if taskTime == nil {
		return &TaskTime{
			T1: time.Now().Day(),
			T2: 20,
		}
	}
	return taskTime
}
