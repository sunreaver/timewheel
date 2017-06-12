# timewheel
Golang时间轮

## 添加

```Golang
PutTimer(second uint, repeat bool, id uint64, e interface{}, callBack CallBackType)

second: second秒之后执行
repeat: 是否循环执行
id: 添加的timer的id
e: callBack的实参
callBack: 执行的方法
```

## 移除

```Golang
RemoveTimer(id uint64)

id: 需要移除执行的timer的id
```
