# douyin 第三届字节青训营项目

---

## 项目介绍：

### 简介：

&emsp;&emsp;**用到的框架和组件:**

<li>微服务框架： go-zero
<li>服务注册中心：Etcd
<li>视频转码： ffmpeg
<li>消息队列： Asynq
<li>存储： Redis, MySQL, Aliyun OSS

&emsp;&emsp;**项目实现的功能如下：**

<li> 视频流接口
<li> 用户注册
<li> 用户登录
<li> 获取用户信息
<li> 投稿接口
<li> 发布列表
<li> 点赞操作
<li> 喜欢列表
<li> 评论操作
<li> 评论列表
<li> 关注操作
<li> 关注列表
<li> 粉丝列表

- 实现上，关注，点赞操作是异步写入mysql的；user信息中的关注数，粉丝数，视频中的点赞数，是定时写入mysql的

&emsp;&emsp;接口文档：https://www.apifox.com/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145
