# robot-2yizi

每隔 5 分钟查询一次录取结果出来了吗? 只支持吉林大学

```go
//实现:
//1. 请求 http://zsb.jlu.edu.cn/Index/gkluqucx.html
//2. 获取一个一次性的 hash, 通过 `meta name="hash"`
//3. 获取验证码图片 base64 之后请求 API 识别, 获取一串数字结果
//4. 请求 http://zsb.jlu.edu.cn/index/gkcxresult.html , 提取响应结果 `.system-message`>
//5. 将结果发送到钉钉的 Webhook, 如果结果没变化在一小时内不会发送重复消息
```

## 如何使用

0. 使用 go 1.16, 安装 app 依赖
   ```shell
   $ cd robot-2yizi
   $ go get
   ```
1. 补全配置文件
   ```go
   const (
       // 考生号
       studentNo string = ``
       // 身份证号
       idNo string = ``
       // 姓名
       name string = ``
   )
   ```
2. 去 [阿里云市场](https://market.aliyun.com/products/57124001/cmapi027426.html) 申请一个验证码识别 API, 可免费使用 30 次, 超额需付费
   ```go
   // 补全 AppCode 配置, 阿里云市场里识别验证码 API 的凭证
   appCode string = ""
   ```
3. 在钉钉群里添加 WebHook, 补全 token
   ```go
    // 钉钉群 Webhook 链接内的 token
    dingToken string = ""
   ```
4. `go run app.go`
   
## 需求:
1. 每次变化的值有两个, 验证码和 Hash
   ```
   yzm	64063510
   hash	1547c50b16cbfc974c778b4a10838236_5cae98a2e901afed5cb12214ea3751e7
   ```

2. hash 需要刷新查询页面(gkluqucx.html), 提交查询后失效
   ```html
   <meta name="hash" content="1547c50b16cbfc974c778b4a10838236_55bdbb7cb9190c41db2ce3449bdce0dc">
   ```
   - hash 的前半段 `1547c50b16cbfc974c778b4a10838236` 是固定的
   - hash 的后半段 `55bdbb7cb9190c41db2ce3449bdce0dc` 直到提交查询前, 刷新多少次都不会再变

3. 验证码可以请求 API 识别, 准确率接近 100%
   - 验证码图片是 http://zsb.jlu.edu.cn/index/verify/tp/1.html 每次请求都会变化

   ```
   <img src="/index/verify/tp/1.html" class="img-rounded" onclick="this.src+='?' + Math.random();" style="cursor:pointer;" title="点击图片刷新验证码">
   ```

## 实现:
1. 请求 http://zsb.jlu.edu.cn/Index/gkluqucx.html
2. 获取一个一次性的 hash, 通过 `meta name="hash"`
3. 获取验证码图片 base64 之后请求 API 识别, 获取一串数字结果
4. 请求 http://zsb.jlu.edu.cn/index/gkcxresult.html , 提取响应结果 `.system-message`>
5. 将结果发送到钉钉的 Webhook.

下面我们先手动进行一次测试, 验证是否能走通流程:
```
yzm=72399438&hash=1547c50b16cbfc974c778b4a10838236_55bdbb7cb9190c41db2ce3449bdce0dc
```
返回
```html
<p class="error">没有查询到你的信息，请核对你的考生号、身份证号、姓名信息是否正确！</p>
```

It works.

## 开发准备:
- 抓包 使用 fiddler
- 从 fiddler 转换到 curl 代码 使用 http://eliaszone.blogspot.com/2016/12/export-to-curl-from-fiddler.html
- 操作 html 使用 https://github.com/PuerkitoBio/goquery
- 图片转 Base64 使用  https://freshman.tech/snippets/go/image-to-base64
- Base64 转图片 使用 https://codebeautify.org/base64-to-image-converter
- curl 转换为 go 代码 使用 https://mholt.github.io/curl-to-go
- api 响应到 go struct 使用 https://mholt.github.io/json-to-go
- json 解析 使用 https://github.com/valyala/fastjson
- 钉钉 webhook 使用 https://github.com/blinkbean/dingtalk
