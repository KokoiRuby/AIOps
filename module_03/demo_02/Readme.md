### ++ Memory

```python
messages = [
        {
            "role": "system",
            "content": "You're an operation & maintenance expert. Your job is to help user solve technical problems."
        },
    ]

while True:
    # blocked, waiting from user input
    user_input = input("Please input your questions:")
    messages.append({
            "role": "user",
            "content": user_input,
        },)
    completion = client.chat.completions.create(
        model="gpt-4-mini",
        messages=messages,
    )

    reply = completion.choices[0].message.content
    print("O&M Expert's answer: ", reply)

    # reply from openai as input for next
    messages.append(
        {
            "role": "assistant",
            "content": reply,
        }
    )
```

