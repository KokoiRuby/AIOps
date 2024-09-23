### Ragflow

[Quickstart](https://ragflow.io/docs/dev/)

```bash
$ git clone https://github.com/infiniflow/ragflow.git
```

```bash
$ cd ragflow/docker
$ chmod +x ./entrypoint.sh

# modify RAGFLOW_VERSION to v0.9.0
$ vi .env

# it shall take a while
$ sudo docker compose up -d

# chk
$ docker logs -f ragflow-server
```

Flow:

1. 注册账号并登录
2. 点击右上角头像，进入模型供应商设置 Model Providers
3. 以 OpenAI 为例，点击添加模型，输入 Base URL 和 API KEY
4. 返回首页，创建知识库 Base Knowledge，选择 Embedding 模型和文档解析方法
5. 上传文件，点击 “启动” 按钮进行 Embedding，等待处理完成
6. 新建助理 Assistant 并基于知识库进行对话