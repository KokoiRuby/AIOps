from openai import OpenAI

client = OpenAI(
    api_key="BLOCKED",
    base_url="https://vip.apiyi.com/v1"
)

completion = client.chat.completions.create(
    model="gpt-4o-mini",
    response_format={
        "type": "json_object",
    },
    messages=[
        {
            "role": "system",
            "content": 'You are a JSON parser. Please output a json object based on example.'
                       'Demo: {"service_name":"","action":""}，'
                       'where action could be get_log, restart, delete'},
        {
            "role": "user",
            "content": "Please help me restart payment service."
        }
    ]
)

print(completion.choices[0].message.content)
