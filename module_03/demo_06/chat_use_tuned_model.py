from openai import OpenAI

client = OpenAI()

# tuned-model
model_id = "ft:gpt-4o-mini-2024-07-18:he3::9tbKC81t"

completion = client.chat.completions.create(
    model=model_id,
    messages=[
        {
            "role": "system",
            "content": "You are a log level identifier, "
                       "please identify the log level for each log entry, output P0, P1 or P2 directly.",
        },
        {"role": "user", "content": "Disk I/O error"},
    ],
)

print("ChatGPT resp: ", completion.choices[0].message.content)
