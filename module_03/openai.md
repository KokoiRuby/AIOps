## [OpenAI](https://platform.openai.com/docs/overview)

### Concepts

- [Text generation models](https://platform.openai.com/docs/concepts/text-generation-models) 文本生成模型 = Generative Pre-trained Transformers 生成式预训练转换器。
- [Assistants](https://platform.openai.com/docs/concepts/assistants) 助理 = 执行用户任务的实体。
- [Embeddings](https://platform.openai.com/docs/concepts/embeddings) 嵌入 = （文本）数据的向量化表示。
- [Tokens](https://platform.openai.com/docs/concepts/tokens) 表示频繁出现的字符序列。
  - [Tokenizer](https://platform.openai.com/tokenizer)

### [Models](https://platform.openai.com/docs/models/continuous-model-upgrades)

### [Libraries](https://platform.openai.com/docs/libraries/python-library)

- [Python](https://github.com/openai/openai-python) (Official)
- [Go](https://github.com/sashabaranov/go-gpt3) (Community)

### Capabilities

#### [Text Generation](https://platform.openai.com/docs/guides/text-generation)

Natural Language/Code/Images understandable.

"Prompts" as input, design of "Prompts" = **Prompt Engineering**.

[Examples](https://platform.openai.com/docs/examples) for inspiration.

temperature ≈ *randomness* (0~2 from most deterministic to least)

#### [Function Calling](https://platform.openai.com/docs/guides/function-calling)

Connect LLM to external tools & systems.

**How to?** [Chat Completions API](https://platform.openai.com/docs/guides/text-generation/chat-completions-api), [Assistants API](https://platform.openai.com/docs/assistants/overview), [Batch API](https://platform.openai.com/docs/guides/batch)

1. 选择代码中需要 LLM 调用的函数
2. 定义函数 JSON Schema
3. JSON Schema 传入 tools 传入 chat completion 通过 create 发起调用
4. 提取响应中的 LLM 返回的参数并发起调用。
5. 将调用结果返回给 LLM

![Function Calling diagram](https://cdn.openai.com/API/docs/images/function-calling-diagram.png)

:pencil2: `strict: true` to enable Structured Outputs follows JSON Schema language. See the [Guide](https://platform.openai.com/docs/guides/structured-outputs/supported-schemas).

:pencil2: `tool_choice` to control calling behavior.

:construction_worker: [Tips and best practices ](https://platform.openai.com/docs/guides/function-calling/tips-and-best-practices)

#### [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs/structured-outputs)

[Avail Models](https://platform.openai.com/docs/guides/structured-outputs/supported-models)

**vs. When to use?**

- 当你希望将 LLM 和工具/函数关联，使用 function calling
- 当你希望结构化模型的输出给用户，使用 `response_format` 指定一个类/类型。

**vs. JSON mode**

结构化输出是 JSON mode 的演进版，仅支持 `4o` series。

**[How to?](https://platform.openai.com/docs/guides/structured-outputs/how-to-use)**

SDK obj (python 推荐) or Manual Schema

:construction_worker: [Tips and best practices](https://platform.openai.com/docs/guides/structured-outputs/best-practices)

#### [Reasoning](https://platform.openai.com/docs/guides/reasoning) (Beta)

`o1` series 支持更复杂的推理，能生成更长的内部 Chain of Thought。

### Endpoints

#### [Chat Completions](https://platform.openai.com/docs/guides/chat-completions)

支持文本 & 图片输入，输出文本内容；通过 `messages` 接受输入。

Message Roles

- `system` 可选，用于设置 assistant 的行为。

- `user` 用户请求 or comment。

- `assistant` 存储上一个 assistant 的响应；用户可以传入自定义内容来控制行为。

  :cry: 每一次 chat 都相当于往 OpenAI 发送一次 HTTP 请求，是无状态的... = LLM 推理每一次也是独立的...

  :smile: `assistant` 可以实现简单的 Memory 功能。

#### [Fine-tuning](https://platform.openai.com/docs/guides/fine-tuning)

> Few-shot Learning = using demonstration to show how to perform a task

:smile: 提供高于 prompting 质量的结果

:smile: 提供更多 examples 用于训练

:smile: Token saving

:smile: 低延迟

[Avail Models](https://platform.openai.com/docs/guides/fine-tuning/which-models-can-be-fine-tuned)

**When to use?** 推荐优先使用 prompt engineering, prompt chaining, function calling.

**How to?**

1. 准备训练数据 `.jsonl`
   - 格式符合 Chat Completins API，具体为 messages list，每条 message 包含 role/content
   - 至少 10 样本；推荐 50 训练样本 on `gpt-4o-mini` & `gpt-3.5-turbo` 决定微调后的模型是否有提升。
   - 需要留意 Token limists & Costs
   - [Format Validation](https://cookbook.openai.com/examples/chat_finetuning_data_prep)
2. 上传训练数据 with `purpose="fine-tune"`
   - 1GB 限制
   - 返回 file ID
3. 创建微调模型
   - 支持 list/retrieve/cancel/list_events/delete
4. 使用微调模型 Chat Completions API
   - 支持 `fine_tuned_model_checkpoint` = 生成模型过程中的版本节点
   - 每一个 checkpoint 都会声明 step_number & metrics

[Analyzing](https://platform.openai.com/docs/guides/fine-tuning/analyzing-your-fine-tuned-model)

#### [Batch](https://platform.openai.com/docs/guides/batch/batch-api)

支持异步分组请求，能有效降低 cost，更高的限速。适用于: 评估，分类大数据集...

[Avail Model](https://platform.openai.com/docs/guides/batch/model-availability)

**How to?**

1. 准备批量文件 `.jsonl`
   - 格式符合 Chat Completins API & Embedding API，具体为 messages list，每条 message 包含 role/content
2. 上传批量文件 with `purpose="batch"`
   - 返回 batch file ID
   - 50,000 requests & 100MB
3. 创建批量任务
   - 可设置完成时间窗口
4. 查看批量任务 [status](https://platform.openai.com/docs/guides/batch/4-checking-the-status-of-a-batch)
5. 获取结构 `.jsonl`

#### Embedding

嵌入 = 测量文本字符串之间的关联性，主要用于

- 搜索（query 关联性排序）
- 集群化（根据相似度分组）
- 推荐（query 关联性）
- 异常检测（最低关联性）
- 多样性测量（相似性分布）
- 分类（label 相似度）

本质：矢量 a list of floating-point numbers；两个矢量的欧式距离大小度量其相似度，越小相似度越高。

Avail Model: `text-embedding-3-small`(length: 1536) & `text-embedding-3-large` (length: 3072)

[Examples](https://platform.openai.com/docs/guides/embeddings/use-cases)