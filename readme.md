# file
## 说明
文件操作框架 用户基本文件io或者状态监控

## file.open
- 打开写的文件句柄 不存在会默认创建
- ud = file.open{name , path , delim}
- 满足lua.writer接口
```lua
    local ud = file.open{
        name = "demo",
        path = "/var/logs/a.log",
        delim = "\n",
    }
    
    start(ud)
```

### 接口函数
- [ud.backup()]()
- [ud.push()]()
```lua
    local ud = file.open("x-yyyy-MM-dd.hh-mm")
    ud.backup() --执行当前目录的事件 根据当前时间
    ud.push("xxx %s" , "helo") --写入
```

## file.stat
- 获取文件状态
- stat = file.stat(path)
### 基础接口
- [stat.ok]()
- [stat.name]()
- [stat.ext]()
- [stat.mtime]()
- [stat.ctime]()
- [stat.atime]()
- [stat.path]()
- [stat.dir]()
```lua
    local st = file.stat("a.txt")
    print(st.ok)
    print(st.name)
    print(st.ext)
    -- etc
```

## file.dir
- 打开文件目录
- d = file.dir(path)
#### 基础接口
- [d.ok]()
- [d.err]()
- [d.count]()
- [d.grep()]()
- [d.ipairs()]()

```lua
    local d = file.dir("/var/logs")
    print(d.err)

    d.grep("*" , function(stat)
        print(stat.ok)        
        --todo
    end)

    d.ipairs(function(stat)
        print(stat.ok)
        --todo
    end)
```

## file.walk
- 新建一个文件遍历器
- walk = file.walk(name)
###基础接口
- [walk.open()](walk.open)

### walk.open
- 打开walk扫描句柄
- tx = walk.open(path)
#### 基础接口
- [tx.ext()]()
- [tx.limit()]()
- [tx.run()]()
```lua
    local wk = file.walk("日志")
    local tx = wk.open("/var/log")
    tx.ext(".log") --匹配指定后缀
    tx.limit(10)   --限速file/s 
    tx.run()       --启动
    
```