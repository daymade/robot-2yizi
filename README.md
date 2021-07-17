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

## 运行效果

```
API server listening at: [::]:14951
hash: 1547c50b16cbfc974c778b4a10838236_d655f758819bd06266b3d716b1131545
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
verificationCode: 39022110
time: 2021-07-17 21:08:04.1554525 +0800 CST m=+6.346174501, message: 没有查询到你的信息
和上次查询结果相同, 不发送消息, 已重复 1 次
hash: 1547c50b16cbfc974c778b4a10838236_94de68aa502e85c0d4dc162c583eac69
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
verificationCode: 29834987
time: 2021-07-17 21:13:16.5909772 +0800 CST m=+318.782139001, message: 没有查询到你的信息
和上次查询结果相同, 不发送消息, 已重复 2 次
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
