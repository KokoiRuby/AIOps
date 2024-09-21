### LangChain Memory

```python
# pip install openai==0.28.1
# pip install langchain==0.0.330

from langchain.chat_models import ChatOpenAI
from langchain.chains import ConversationChain
from langchain.memory import ConversationBufferMemory，ConversationBufferWindowMemory

# build llm
llm = ChatOpenAI(
    model="gpt-4o-mini",
    openai_api_key="BLOCKED",
    openai_api_base="https://vip.apiyi.com/v1",
    max_tokens=None,
    timeout=None,
)

# create memory
memory = ConversationBufferMemory()
# keep latest 30 chats
memory = ConversationBufferWindowMemory(k=30)
# limit token cnt
memory = ConversationTokenBufferMemory(llm=llm, max_token_limit=128000)
# summarize if over limit
memory = ConversationSummaryMemory(llm=llm, max_token_limit=128000)

# create chain, inject llm & memory
chain = ConversationChain(
    llm=llm,
    memory=memory,
    verbose=False,
)

while True:
    user_input = input("Please input your questions: ")
    reply = chain.predict(input=user_input)
    print("Current context is: ", chain.memory)

```

