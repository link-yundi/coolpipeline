package coolpipeline

import (
	"fmt"
	"github.com/link-yundi/ylog"
	"testing"
	"time"
)

/**
------------------------------------------------
Created on 2022-11-04 18:50
@Author: ZhangYundi
@Email: yundi.xxii@outlook.com
------------------------------------------------
**/

func TestName(t *testing.T) {
	addFunc := func(in any) (out any) {
		i := in.(int)
		tmp := i
		i++
		ylog.Info(tmp, "add", 1, "=", i)
		return i
	}
	squareFunc := func(in any) (out any) {
		i := in.(int)
		tmp := i
		i *= i
		ylog.Info(tmp, "square", "=", i)
		return i
	}
	// 工作流定义
	pool := NewPipelines(1, addFunc, squareFunc, addFunc)
	pool.AddTask(1, 3, 6)
	pool.Wait()
}

func TestFactor(t *testing.T) {
	start := time.Now()
	// 采购
	buy := func(in any) (out any) {
		time.Sleep(1 * time.Second)
		i := in.(int)
		out = fmt.Sprint("零件", i)
		ylog.Info(out)
		return
	}
	// 组装
	build := func(in any) (out any) {
		time.Sleep(5 * time.Second)
		out = "组装(" + in.(string) + ")"
		ylog.Info(out)
		return
	}
	// 打包
	pack := func(in any) (out any) {
		time.Sleep(3 * time.Second)
		out = "打包(" + in.(string) + ")"
		ylog.Info(out)
		return
	}
	// 工作流定义
	pipeline := NewPipelines(2, buy, build, pack)
	// 订购3台
	var ins []any
	for i := 1; i <= 3; i++ {
		ins = append(ins, i)
	}
	pipeline.AddTask(ins...)
	pipeline.Wait()
	end := time.Now()
	duration := end.Sub(start).Seconds()
	ylog.Info("耗时: ", duration, "s")
}
