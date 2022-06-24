# cond
表达式快速匹配模块 解决 过于冗长的条件判断 简化 开发逻辑

## vela.cond(string , string , string)
- 编译一个条件表示的功能模块
- cond 实现了 NewIndex的方法 支持直接赋值修改
- ok 匹配成功后触发相关逻辑 类型:pipe 参数:match传入的参数
- no 没又匹配后触发相关逻辑 类型:pipe 参数:match传入的参数
- match(object) 匹配

```lua
    local cnd = vela.cond("name eq zhang1,zhang2,zhang3")
    
    cnd.ok(function(v) print(v.name , v.pass) end)
    cnd.no(function(v) print(v.name , v.pass) end)
    
    local object = {
        name = "wangwu",
        pass = "123456",
    }
    cnd.match(object)
```

### 条件格式
- key,key,key 键值采用逗号分割 代表多个
- val,val,val 匹配那些值也可以采用逗号 代表多个
- 中间 为 表达是 缩写 包含 eq,re,cn,in,lt,le,ge,gt
- eq 等于规则(aa eq aa,bb,cc)
- re 匹配规则(abc re a*,*b*,*c) 只有星号\*和\.b 不支持正则表达式
- cn 包含内容(abc cn *c)
- in 等于包含(1 in 1,3,4)
- lt 小于
- le 小于等于
- ge 大于等于
- gt 大于
```
    # 单个
    key eq value 
    key eq value,value,value
    
    # 多个
    key1,key2 eq value1,value2
    
    # 取反
    key !eq value,value2,value3
    key !cn *v,*v2
```