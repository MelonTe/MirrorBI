<h1 style="text-align: center; white-space: nowrap;">
  <img
    src="README.assets/MBI.ico"
    alt="Logo"
    width="60"
    height="60"
    style="display: inline-block; vertical-align: middle; margin-right: 8px;"
  />
  <span style="display: inline-block; vertical-align: middle;">MirrorBI</span>
</h1>

## 1、项目介绍

MirrorBI 是一个智能数据分析平台，面向非技术用户提供“一键式数据洞察”能力：用户仅需上传原始 Excel/CSV 数据集并输入业务目标（Goal），系统即可自动完成——1.Prompt 生成 → 2. LLM 数据分析 → 3. ECharts 图表渲染 → 4. 文字结论输出——全过程，实现零 SQL、零建模的自助式 BI。

前后端已经打包，具体查看代码详情。后端需要编辑Config.yaml文件，否则将无法正常运行。模板如下：

```yaml
#数据库配置
database:
  user: "xxx"
  password: "xxx"
  host: "xxx"
  port: 3306
  name: "mirrorbi"
rds:
  host: "xxx"
  port: xxx
  username: "xxx"
  password: "xxx"
tcos:
  bucketName: "xxx"
  region: "xxx"
  host: "https://xxxxxxx.myqcloud.com"
Siliconflow:
  APIkey: "xxxx"
RabbitMQ:
  host: "xxx"
  port: 5672
  username: "xxx"
  password: "xxx"
```

## 2、项目整体功能展示

### 2.1、样本数据准备

准备了两个样本数据，分别是日期-人数Excel表，以及模拟出的实习情况分析Excel表。

![image-20250519164042711](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164042711.png)

![image-20250519164051922](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164051922.png)

### 2.2、智能分析界面（同步）

界面样式如下：

![image-20250519164147450](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164147450.png)

输入分析目标、图表名称、以及图表类型和数据源。

![image-20250519164221641](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164221641.png)

然后点击“**开始智能分析按钮**”，数据和需求会在后端**自动解析**，**嵌入Prompt发送给LLM模型，进行处理**。期间需要**同步等待响应**。

![image-20250519164335259](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164335259.png)

等待一定时间后，会跳出“分析成功”字段，然后页面自动渲染出分析结果。

![image-20250519164508511](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164508511.png)

### 2.3、智能分析界面（异步）

除了同步等待响应结果，还可以选择异步分析界面，将数据提交后，可以直接提交下一个任务，做到不阻塞等待，优化用户体验。

界面样式与同步的相同。

![image-20250519164708495](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164708495.png)

点击分析按钮后，会直接跳转到我的图表界面，此时**可以观察到正在等待执行或者正在执行的任务。**

![image-20250519164752094](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164752094.png)

### 2.4、我的图表界面

#### 2.4.1、图表执行状态

在异步提交任务后，可以看到正在被处理的任务，以及任务的执行情况。

![image-20250519164846068](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519164846068.png)

点击**“刷新状态”**按钮，可以**实时加载任务执行情况**。最右侧有“删除”按钮图表，可以选择性删除任务。

#### 2.4.2、历史图表查询

在执行状态的下边，是所有历史图表的生成情况，可以进行简单的预览。

![image-20250519165028256](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519165028256.png)

![image-20250519165033979](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519165033979.png)

### 2.5、图表详情页查看

对选中的图表卡片，点击**“展示”**按钮，可以进入到图表AI分析详情页查看具体的分析信息。

![image-20250519165143768](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519165143768.png)

点击**“展示原始数据”**按钮，可以查看上传的Excel解析后的数据。

![image-20250519165309400](C:\Users\minat\Desktop\Note\MirrorBI\README.assets\image-20250519165309400.png)

