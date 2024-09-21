from langchain.chat_models import ChatOpenAI
from langchain.chains import ConversationChain
from langchain.memory import ConversationBufferMemory

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
