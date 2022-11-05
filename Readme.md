# CoolPipeline

流水线：并发、池化

### 安装

```sh
go get -u github.com/link-yundi/coolpipeline
```

### 特性

- 池化，复用，减少GC压力
- 流水线定制
- 多条流水线并发

### 示例

#### 工作流：

```go
import "github.com/link-yundi/ylog"

func main() {
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
    // 流水线定义： 工作顺序
    pool := NewPool(1, addFunc, squareFunc, addFunc)
    pool.AddTask(1, 3, 6)
    pool.Wait()
    pool.Close()
}
// output:
// 2022/11/05 13:49:41 [INFO] [1 add 1 = 2]
// 2022/11/05 13:49:41 [INFO] [3 add 1 = 4]
// 2022/11/05 13:49:41 [INFO] [4 square = 16]
// 2022/11/05 13:49:41 [INFO] [16 add 1 = 17]
// 2022/11/05 13:49:41 [INFO] [6 add 1 = 7]
// 2022/11/05 13:49:41 [INFO] [7 square = 49]
// 2022/11/05 13:49:41 [INFO] [49 add 1 = 50]
// 2022/11/05 13:49:41 [INFO] [2 square = 4]
// 2022/11/05 13:49:41 [INFO] [4 add 1 = 5]
```

#### 并发：

```go
import "github.com/link-yundi/ylog"

// 定制 3 台手机, 两套流水线
// 步骤： 采购零件(耗时1s) -> 组装(耗时5s) -> 打包(耗时3s), 一套完整的流程耗时 9s

func main() {
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
    // 工作流定义: 2个并发
    pool := NewPool(5, buy, build, pack)
    // 订购3台
    var ins []any
    for i := 1; i <= 3; i++ {
        ins = append(ins, i)
    }
    pool.AddTask(ins...)
    pool.Wait()
    end := time.Now()
    duration := end.Sub(start).Seconds()
    ylog.Info("耗时: ", duration, "s")
}

// output
// 2022/11/05 13:55:23 [INFO] [零件2]
// 2022/11/05 13:55:23 [INFO] [零件1]
// 2022/11/05 13:55:24 [INFO] [零件3]
// 2022/11/05 13:55:28 [INFO] [组装(零件1)]
// 2022/11/05 13:55:28 [INFO] [组装(零件2)]
// 2022/11/05 13:55:29 [INFO] [组装(零件3)]
// 2022/11/05 13:55:31 [INFO] [打包(组装(零件2))]
// 2022/11/05 13:55:31 [INFO] [打包(组装(零件1))]
// 2022/11/05 13:55:32 [INFO] [打包(组装(零件3))]
// 2022/11/05 13:55:32 [INFO] [耗时:  10.004996625 s]
```

