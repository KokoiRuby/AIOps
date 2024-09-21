## [Doc](https://platform.openai.com/docs/overview)

[API key](https://api.apiyi.com/) (Domestic)

### Chat Completions

Role 角色

- system 系统设定 Prompt
- user 发起对话的用户
- assistant 基于用户消息模型返回的响应

:confused: 如何基于 AI 的回答继续对话？（比如我认为回答不对...）****

- 每一次 chat 都相当于往 OpenAI 发送一次 HTTP 请求，是无状态的... = LLM 推理每一次也是独立的...

:smile: ++ Memory 即将上一次 LLM 的输出作为下一次的输入添加到 messages 中，标记 role 为 assistant 并再次发起一次 chat。

#### [LangChain Memory](https://python.langchain.com/v0.1/docs/modules/memory/)

Python SDK provides utilities for adding memory to LLM.

- `ConversationBufferMemory` 最简单，通过传递历史消息来实现记忆，程序退出后对话上下文不会被保存
- `ConversationBufferWindowMemory` 根据对话次数控制记忆长度
- `ConversationTokenBufferMemory` 根据 Token 控制记忆长度
- `ConversationSummaryMemory` 对超限 token 进行总结（额外开销）
- ，实现更长的记忆

### JSON Mode

LLM 稳定输出 JSON 格式

- `client.chat.completions.create` ++ `response_format={"type": "json_object"}`
- system prompt 须包含 JSON 关键字

:smile: 从语义化内容中提取参数，然后交由程序处理；适合 few-shot。

### Function Calling

:cry: JSON Schema 相对固定，难以实现复杂业务逻辑。对于每一种业务可能都要编写一个 schema...

:smile: LLM 选择合适的函数并自动填充参数，将 LLM 和程序关联起来，在 tools 中定义函数（数组），并传入 chat completion 中。

**注：LLM 返回的只是填充后的函数调用，程序中需要解析并发起调用。**即解析 LLM 返回中的函数名和参数

### Fine-tuning

微调：通过给大模型更多的样本数据来改进模型的表现，使其在一些特定的任务上输出更好的结果。

- 要求输出特定的风格和格式
- 提高输出的可靠性
- 纠正无法遵循复杂提示语的问题
- 执行一项难以用语言表达的新技能或任务

Flow:

1. **准备微雕数据，即 `jsonl` 包含数条需（转换成） fine-tuning 的数据**
2. 将 `jsonl` 上传到 OpenAI
3. 创建一个 Job 返回一个 Job ID
4. 根据 Job ID 查询状态，最终会返回给我们一个微调后 Model

可微调模型：`gpt-4o-mini-2024-07-18` & `gpt-3.5-turbo-0125`。