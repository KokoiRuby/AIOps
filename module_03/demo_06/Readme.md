### Fine-tuning

根据真实的已标注日志数据，使模型具备判断日志故障优先级的能力。

```csv
log,priority
[2024-08-07 12:00:00] Database connection failed. System is down.,P0
[2024-08-07 12:05:00] Unrecoverable error in payment processing module. System halted.,P0
[2024-08-07 12:10:00] Authentication service is not responding. All user access is blocked.,P0
...
```

```python
import csv
import json
import os

def convert_csv_to_jsonl(csv_file_path, jsonl_file_path):
    with open(csv_file_path, mode="r", encoding="utf-8") as csv_file:
        csv_reader = csv.DictReader(csv_file)

        # for each row build an entry represents a message to LLM
        with open(jsonl_file_path, mode="w", encoding="utf-8") as jsonl_file:
            for row in csv_reader:
                log = row["log"]
                priority = row["priority"]

                jsonl_entry = {
                    "messages": [
                        {
                            "role": "system",
                            "content": "You are a log level identifier, "
                                       "please identify the log level f each log entry, output P0, P1 or P2 directly.",
                        },
                        {"role": "user", "content": log},
                        {"role": "assistant", "content": priority},
                    ]
                }

                # append to row jsonl
                jsonl_file.write(json.dumps(jsonl_entry, ensure_ascii=False) + "\n")


csv_file_path = os.path.join(os.path.dirname(__file__), "./data", "log.csv")
jsonl_file_path = "./log.jsonl"
convert_csv_to_jsonl(csv_file_path, jsonl_file_path)
```

```python
from openai import OpenAI

client = OpenAI(
    api_key="BLOCKED",
    base_url="https://vip.apiyi.com/v1"
)

# upload jsonl to llm with purpose
file_name = client.files.create(file=open("log.jsonl", "rb"), purpose="fine-tune")
file_id = file_name.id
print("file_id is: ", file_id)

# create fine-tune job
finetune_job = client.fine_tuning.jobs.create(
    training_file=file_id, model="gpt-4o-mini"
)
job_id = finetune_job.id
print("job_id is: ", job_id)
```

```python
from openai import OpenAI

client = OpenAI()

job_id = "ftjob-q8AhUAUIzdlaeLAPlU3yzcmv"
status = client.fine_tuning.jobs.retrieve(job_id)

print("\n", status)
print("\nfinetune status is: ", status.status)
print("\nfinetune model id is: ", status.fine_tuned_model)
```



