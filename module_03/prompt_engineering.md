### Prompt

**Prompt 提示语**：给人工智能模型输入文本或指令，引导模型生成特定的输出。

:smile: Prompt = A good question = A good ctxt = A good answer 即一个优秀的 prompt 决定了模型输出的质量。

**Flow: 基于目标和例子，根据用户输入提供输出；**

:notebook: [Wonderful-prompt](https://github.com/langgptai/wonderful-prompts)

### Token

当编写好 Prompt 之后，以 Token 的形式提交给 LLM；

LLM 计算 Token 中的统计关系，并基于已输出的 Token 计算下一个 Token 的概率分布，选择概率最高的 Token 作为预测结果。

**注：空格也会占用 Token。**

#### NPL

NLP 自然语言处理时 AI 的一个分支，目的是让计算机理解生成并处理人类语言。

Tokenization 分词是将一串文本进行分解成 token（若干单词的组合）。

NLP 中最简单的分词算法，即**按空格分词** "I am a student" → ['I', 'am', 'a', 'student']

:smile: 保留词的语义。

:cry: 词表随训练集大小线性增大；若遇到不在词表中的字符时会出现 OOV out-of-vocabulary 导致无法正确处理。

NLP 还可以**按单字符进行分词**，"I am a student" → ['I', 'a', 'm' ...]

:cry: 不保留语义

#### BPE

:cry: 若把每个单词视为一个 token 对应一个整数，那么随着不同国家的语言加入和词汇增加，词表会变得越来越大。

:smile: 字节对编码 Byte Pair Encoding 能够高效减少 token 数量。

1. 构建词表

```bash
# 确定词表大小，即单词/子单词的个数
# 每个单词以 </w> 结尾，统计出现的次数
# 将每个单词拆分为单个字符
there is always a better way
# ↓
t h e r e </w>   1
i s </w>         1
a l w a y s </w> 1 
a </w>           1
b e t t e r </w> 1 
w a y </w>       1 
```

```bash
# 挑出频次最高的字符对，将词表中所有符合的字符对合并成一个子单词

# 第 1 次迭代
# 频次最高字符对 e & r 组合成一个子单词 er，更新词表
t h er e </w>    1 
i s </w>         1
a l w a y s </w> 1 
a </w>           1
b e t t er </w>  1 
w a y </w>       1
 
# 第 2 次迭代
# 频次最高字符对 s & </w> 组合成一个子单词 s</w>，更新词表
t h er e </w>   1
i s</w>         1
a l w a y s</w> 1 
a </w>          1
b e t t er </w> 1 
w a y </w>      1

# 第 3 次迭代
# 频次最高字符对 w & a 组合成一个子单词 wa，更新词表
t h er e </w>   1
i s</w>         1
a l wa y s</w>  1 
a </w>          1
b e t t er </w> 1 
wa y </w>       1 

# 第 4 次迭代
# 频次最高字符对 wa & y,然后更新词表为:
t h er e </w>   1
i s</w>         1
a l way s</w>   1 
a </w>          1
b e t t er </w> 1 
way </w>        1 

# 直到词表中数量达到阈值 V 或者下一个字符对的频数为 1
```

2. 编码

```bash
# 1. 将词表中的单词按长度大小进行排序
# 2. 对于待编码的单词，遍历排序好的词表，判断词表中的单词/子词是否是待编码单词的子串，如果是，则输出当前子词，并继续遍历词表。
# 3. 如果遍历完词表，单词中仍然有子字符串没有被匹配，那我们将其替换为一个特殊的子词，比如 <nok>

# 假设当前构建好的词表为：
er</w> 1
be     1
t      1
 
better</w>
# ↓
be t t er</w>
```

3. 解码

```bash
# 将所有的输出子词拼在一起, </w> 视为空格
th er e<w> is
# ↓
there is
```

#### [tiktoken](https://github.com/openai/tiktoken)

A fast [BPE](https://en.wikipedia.org/wiki/Byte_pair_encoding) tokeniser for use with OpenAI's models.

#### Context Window

上下文窗口：Token 输入/输出是有限的，LLM 不可能记忆所有的内容。

窗口越大，可输入的内容以及可容纳的 Token 就越多。

**注：模型回答效果并不会随着输入 Token 数量线性而提升。**

### [Prompt Engineering](https://platform.openai.com/docs/guides/prompt-engineering)

[Write clear instructions](https://platform.openai.com/docs/guides/prompt-engineering/write-clear-instructions) 明确告知模型应该怎么做

- 包含尽可能多地细节，避免让 Model 进行猜想
- 让模型具备一个人格
- 使用分隔符比如 """ XML <tag> 界定 Model 的输出内容
- 指明需要完成一个任务的步骤 = Chain of Thoughts
- 提供例子
- 指明输出的长度

[Provide reference text](https://platform.openai.com/docs/guides/prompt-engineering/provide-reference-text) 提供参考性的文本

- 告诉 Model 借助参考文档来回答
- 告诉 Model 使用参考文档中的引用来回答

[Split complex tasks into simpler subtasks](https://platform.openai.com/docs/guides/prompt-engineering/split-complex-tasks-into-simpler-subtasks) 分解复杂任务成简单子任务

- 对用户的 query 使用带意图的分类来标识最相关的指令
- 对于长对话的应用，总结 or 过滤之前的对话内容
- 分片总结长文档，递归构建总结

[Give the model time to "think"](https://platform.openai.com/docs/guides/prompt-engineering/give-the-model-time-to-think) 给予模型足够的时间思考

- 先让 Model 判断 solution 是否正确，再去计算
- 使用内独白 or 查询序列来隐藏模型的推理过程 step1/2/3...
- 询问 Model 是否遗漏了之前的内容

[Use external tools](https://platform.openai.com/docs/guides/prompt-engineering/use-external-tools) 使用外部工具

- 使用基于 Embedding 的搜索实现高效的知识获取
- 使用代码执行去实施更准确的计算 or 外部 API 调用
- 让 Model 访问具体的函数 function calling

[Test changes systematically](https://platform.openai.com/docs/guides/prompt-engineering/test-changes-systematically) 系统性测试变化

- [OpenAI Evals](https://github.com/openai/evals)
- 参加黄金标准答案评估 Model 输出