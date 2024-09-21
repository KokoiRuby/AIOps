### Chat Completions

[Quickstart](https://platform.openai.com/docs/quickstart)

```python
from openai import OpenAI

# client to OpenAI
client = OpenAI(
    api_key="BLOCKED",
    base_url="https://vip.apiyi.com/v1"
)

# create a chat completion
completion = client.chat.completions.create(
    # model to use
    model="gpt-4o-mini",
    # A list of messages comprising the conversation so far
    messages=[
        {
            "role": "system", 
            "content": "You're an operation & maintenance expert. Your job is to help user solve technical problems."},
        {
            "role": "user",
            "content": "Linux ports are in conflict, how to troubleshoot? Just give me a few commands."
        }
    ]
)

print(completion.choices[0].message.content)
```

