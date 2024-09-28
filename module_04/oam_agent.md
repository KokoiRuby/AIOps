## Impl

1. **Self-hosting**: QAnything, RAGFlow, Open Web UI

2. **SaaS**: [Coze](https://www.coze.com/space/7419275647354732560/knowledge), Dify

3. **Secondary Dev** by integrating LangChain & LangGraph

   a. Create a ticket

   b. Read internal OAM knowledge base

   c. Open IM tool

   d. Cause anlysis

### Coze

1. 注册 https://www.coze.com/
2. 打开 Personal -> Knowledge -> Create Knowledge 上传知识库文件并进行 Embedding
3. 创建 Bot
4. 在 Knowledge -> Text 下点击 + ，添加刚才创建的知识库
5. 在 Source 栏开启 Show source
6. 在右侧开始调试，提问：你知道 payment 服务是谁维护吗？
7. 发布 Bot，获得一个公网访问链接，例如：https://www.coze.com/store/bot/7404669938810142727?bot_id=true

**Knowledge settings**

- **Call method**
  - `Auto-call` 自动调用 RAG（可能不精准）
  - `On-demand call` 按需，在提示语里描述什么时候使用知识库
- **Search strategy**
  - `Hybird` 结合全文检索和语义检索的优势，并对结果进行综合排序召回相关的内容片段
  - `Semantics` 像人类一样去理解词与词，句与句之间的关系
  - `Full Text` 基于关键词进行全文检索
- **Maximum recalls** 选择从检索结果中返回多少个内容片段给大模型使用
- **Minimum matching degree** 根据设置的匹配度/相似度选取要返回给大模型的内容片段

