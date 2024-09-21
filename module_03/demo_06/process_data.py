import csv
import json
import os

def convert_csv_to_jsonl(csv_file_path, jsonl_file_path):
    with open(csv_file_path, mode="r", encoding="utf-8") as csv_file:
        csv_reader = csv.DictReader(csv_file)

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

                jsonl_file.write(json.dumps(jsonl_entry, ensure_ascii=False) + "\n")


csv_file_path = os.path.join(os.path.dirname(__file__), "./data", "log.csv")
jsonl_file_path = "./log.jsonl"
convert_csv_to_jsonl(csv_file_path, jsonl_file_path)
