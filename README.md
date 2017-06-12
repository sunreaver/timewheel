# timewheel
Golang时间轮

## 回调方法

```Golang
type CallBackType func(e interface{})
```

## 添加

```Golang
PutTimer(second uint, repeat bool, id uint64, e interface{}, callBack CallBackType)
```

参数 | 意义
--- | ---
**second:** | _second秒之后执行_
**repeat:** | _是否循环执行_
**id:** | _添加的timer的id_
**e:** | _callBack的实参_
**callBack:** | _执行的方法_

## 移除

```Golang
RemoveTimer(id uint64)
```

**id:** _需要移除执行的timer的id_
