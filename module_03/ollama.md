### [Ollama](https://github.com/ollama/ollama)

[Models](https://ollama.com/library) 注：7B 模型至少需要 8G RAM、13B 至少需要 16G RAM、33B 至少需要 32G RAM。

Ollama 支持接入接入任何开源 GGUF 模型 ← :hugs: [Hugging Face](https://huggingface.co/)

```bash
# create Modelfile and specify local GGUF model
FROM ./vicuna-33b.Q4_0.gguf

# create local model
$ ollama create example -f Modelfile

# run model
$ ollama run example
```

```bash
# ali qwen
$ ollama run qwen2:7b

# deepseek
$ ollama run deepseek-coder-v2
```

#### Customize System Prompt

```bash
# 拉取模型
$ ollama pull llama3.1

# create Modelfile & specify model & prompt

# create model by Modelfile
$ ollama create mario -f ./Modelfile

# run model
$ ollama run mario
```

```dockerfile
FROM llama3.1

PARAMETER temperature 1

SYSTEM """
You are Mario from Super Mario Bros. Answer as Mario, the assistant, only.
"""
```

#### REST API

兼容 OpenAPI 支持一次性和对话推理。

```bash
# one-off
$ curl http://localhost:11434/api/generate -d '{
	"model": "llama3.1",
	"prompt": "Why is the sky blue?"
}'

# chat
$ curl http://localhost:11434/api/chat -d '{
	"model": "llama3.1",
	"message": [
		{
			"role": "user",
			"content": "why is the sky blue?"
		}
	]
}'
```

#### [Open WebUI](https://github.com/open-webui/open-webui)

http://localhst:3000/

支持 RAG = Ollama + Qwen2 + Open WebUI

支持 **#网页** 引用网页内容

More

- 支持搜索引擎，例如 Google
- 支持配置 Ollama 多实例负载均衡
- 支持 TTS 和语音识别
- 支持集成 ComfyUI 实现文生图
- 支持多模型对话

```bash
# if Ollama is on your computer
$ docker run -d -p 3000:8080 \
	--add-host=host.docker.internal:host-gateway \
	-v open-webui:/app/backend/data \
	--name open-webui \
	--restart always \
	ghcr.io/open-webui/open-webui:main
```





