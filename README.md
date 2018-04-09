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
second | 间隔
repeat | 是否循环执行
id | 添加的timer的id（用于移除）
e | callBack的实参
callBack | 执行的方法

## 移除

```Golang
RemoveTimer(id uint64)
```

参数 | 意义
--- | ---
id | 需要移除执行的timer的id
