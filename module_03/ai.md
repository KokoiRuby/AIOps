## Generative AI

### Overview

AI refers to develop a computer system which could perform tasks that typically require human intelligence.

Subfields:

- **Machine Learning 机器学习**，如何和让计算机能够在没有明确编程的情况下从数据中学习 & 改进
  - Deep Learning 深度学习，使用人工神经网络
- **Natural Language Processing 自然语言处理**，如何让计算机理解/解释/生成人类语言。
  - Language Generation 语言生成
  - Answering Questions 问题回答
  - Text Classification 文本分类
  - Sentiment Analysis 情感分析
  - Machine Translation 机器翻译
  - Deep Learning Model 深度学习模型，尤其是 **Transformer 模型 (Encoder + Decoder)** 给 NLP 带来了革命性的变化
- **Computer Vision & Audio Processing 计算机视觉 & 音频处理**

### GenAI

**生成式 AI** 即以内容（文本/图像/音视频）生成为主的 AI

**LLM Large Language Model 大语言模型/大模型**：基于大量文本数据训练的人工智能模型，用于 QA/Coding...

**Prompt Engineering 提示语工程**，以便从 GenAI 模型从中获取准确的输出。

### API

大多数 GenAI 模型都可以通过获取对应平台的 API KEY，使用 REST API 访问。

- 妥善保管 KEY
- Tradeoff Token & Rate limits & Cost

### DIY

**RAG 检索增强生成**：提供额外的 Knowledge Base 知识库到 System Prompt 中改进模型推理结果。Like an expert, he/she knows a lot things & he/she knows where to get'em.

**Fine Tuning 微调**：提供具有能够处理特定样本数据的模型。

## Tokenization

LLM 中的一切都是以 Token 形式存在，即模型理解输入 & 处理的基本单元，可以是文本/图片/音频/视频...。

每个 Token 被赋予数值或标识符，并按序列或向量排列。

### Text

**Tokenization 分词**即 Text → Token 的过程，将文本分割成更小的单元。

1. **预处理**：去噪（特殊符号标点）+ 标准化（大小写转换）
2. **分词**：分割文本成 Token
3. 标记：标识 Token

经验：一个 token 相当于一个英文单词的 3/4

### Picture

Picture → Token 需要借助计算机视觉模型

1. **预处理**：对图像进行预处理，包括调整大小、裁剪和归一化
2. **特征提取**：使用 CNN (卷积神经网络) 提取图像的特征，每个特征对应一个 Token

### Video

Video → Token 类图片转换，但需要处理连续帧和实践序列。专用 Tokenizer MAGVIT-v2

1. **帧提取**：从视频中提取每一帧。
2. **帧转换**：将每一帧转换为 Token，类似于图片转换过程中的特征提取和嵌入。
3. **序列化**：将所有帧的 Token 序列化，形成一个连续的 Token 序列，以便进行进一步处理。

## [Transformer](https://explainer.tubex.chat/)

Transformer 是一种彻底改变人工智能方法的神经网络架构，于 2017 "Attention is All You Need" 中提出，并称为深度学习的首选框架。

原理：生成内容取决于对下一个词的预测 Predict。

组件：

1. **Embedding**：Text → Token → Vector（带语义）
2. **Transformer Block**：模型处理和转换输入的基本构建块
   - Multi-head Attention 多头注意力：Transformer 核心模块，允许 Token 间通信，捕获上下文关系
   - Multi-layer Perception 多层感知：前馈网络，对每个 Token 独立操作，细化表示。
3. **Probability**: Linear 层 + Softmax 层 将处理过 Embedding 转换成概率，使模型能够对序列中的下一个 Token 进行预测。

### Embedding

1. Tokenization 分词构建词汇表，每一个带一个 ID
2. Token Embedding 每一个 Token 表示为一个 X 维度的向量，所有 Token 保存在 **(总 Token 数 x Dimension of Vector) 矩阵**中。
3. Positional Encoding 编码位置信息
4. Embedding 最终表现形式

### Transformer Block

#### Multi-head Attention

多头注意力：关注输入序列的相关部分，从而捕获数据中的复杂关系和依赖关系。

1. QKV：每个 token 的 Embedding 向量（通过 Embedding 矩阵分别和 QKV 学习矩阵相乘）被转换为 QKV 三个向量。

- `Query` 表示当前想要寻找关联上下文的词。
- `Key` 表示序列中所有词的上下文。
- `Value` 包含每个词关联的真实信息和具体表示。

2. Masked Self Attention 掩码自注意力。

- `Attention Score` Q 和 K 矩阵的点积决定了每个 Q 与 K 的对齐，产生一个正方形矩阵，反映了所有输入 Token 之间的关系。
- `Mask` 将模型无法访问的未来 token 的值设置为负无穷大，模型学习如何在不 "偷看" 未来的情况下预测下一个 Token。
- `Softmax` 将 `Attention Score` 转换成概率，即取指数，矩阵的每一行总和为 1，并指示每个其他 token 对其左侧每个 token 的相关性。

3. Output 输出

- `Attention Score` x V 矩阵

#### Multi-layer Perception

多层感知：将自注意力表示投影到更高维度，以增强模型的表示能力。由多个线性变换组成，将维度先增加减少。

MLP 独立地处理 Token，并简单地将它们从一个表示映射到另一个表示。

### Probability

输入经过所有 Transformer 块处理后，输出通过最后的线性层以准备 token 预测。

该层将最终表示投影到一个 X 维的空间中，每一个 Token 对应一个 logit 概率。

应用 logit 于 softmax 转换成概率分布，总和为 1。logit 决定了每个 token 成为序列中下一个词的可能性。

**temperature**

- `=1` 无影响。
- `>1` 较低的temperature使模型更加自信和确定，通过锐化概率分布来导致更可预测的输出。
- `<1` 较高的temperature创建一个更软的概率分布，允许生成的文本中有更多的随机性和多样性。