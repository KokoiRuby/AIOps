from openai import OpenAI

client = OpenAI(
    api_key="BLOCKED",
    base_url="https://vip.apiyi.com/v1"
)

messages = [
        {
            "role": "system",
            "content": "You're an operation & maintenance expert. Your job is to help user solve technical problems."
        },
    ]

while True:
    user_input = input("Please input your questions: ")
    messages.append({
            "role": "user",
            "content": user_input,
        },)
    completion = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=messages,
    )

    reply = completion.choices[0].message.content
    print("O&M Expert's answer: ", reply)

    messages.append(
        {
            "role": "assistant",
            "content": reply,
        }
    )
