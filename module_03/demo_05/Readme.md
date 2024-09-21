### Function Calling

```python
from openai import OpenAI

client = OpenAI(
    api_key="BLOCKED",
    base_url="https://vip.apiyi.com/v1"
)

# input
query = 'Check logs contains keyword "Error" and labeled with app=grafana.'
query = 'Check logs contains keyword "Fail" and labeled with app=payment.'

messages = [
        {
            "role": "system",
            "content": "You're a Loki log querier, u could help user to analyze log, "
                       "u could call multiple functions to finish tasks from user, and try to analyze the cause."},
        {
            "role": "user",
            "content": query
        }
    ]

tools = [
    {
        "type": "function",
        "function": {
            "name": "analyze_loki_log",
            "description": "Get logs from Loki",
            "parameters": {
                "type": "object",
                "properties": {
                    "query_str": {
                        "type": "string",
                        "description": 'Loki query string, example: {app="grafana"} |= "Error"',
                    },
                },
                "required": ["query_str"],
            },
        },
    }
]

response = client.chat.completions.create(
    model="gpt-4o-mini",
    messages=messages,
    tools=tools,
    tool_choice="auto",
)
response_message = response.choices[0].message
tool_calls = response_message.tool_calls

print("\nChatGPT want to call function: ", tool_calls)

```

