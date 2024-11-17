import asyncio
import aiohttp
import random
import schedule
import time

# Corpus of clean random sentences
corpus = [
    "The quick brown fox jumps over the lazy dog.",
    "A journey of a thousand miles begins with a single step.",
    "To be or not to be, that is the question.",
    "All that glitters is not gold.",
    "An apple a day keeps the doctor away.",
    "Better late than never, but never late is better.",
    "Actions speak louder than words.",
    "Do not count your chickens before they hatch.",
    "Honesty is the best policy.",
    "When in doubt, leave it out.",
    "Every cloud has a silver lining.",
    "Necessity is the mother of invention.",
    "Fortune favors the bold.",
    "You reap what you sow.",
    "A picture is worth a thousand words.",
    "Brevity is the soul of wit.",
    "Practice makes perfect.",
    "Rome wasn't built in a day.",
    "Still waters run deep.",
    "Time waits for no one.",
]


# Function to select a random sentence from the corpus
def random_sentence():
    return random.choice(corpus)


# Function to send a POST request to the backend
async def send_message():
    print("Attempting to send a message")
    url = "http://go-server:8080/send-message"  # Backend service endpoint
    message = {"text": random_sentence()}
    async with aiohttp.ClientSession() as session:
        try:
            async with session.post(url, json=message) as response:
                if response.status == 200:
                    print(f"Message sent successfully: {message['text']}")
                else:
                    print(f"Failed to send message: {response.status}")
        except Exception as e:
            print(f"Error sending message: {e}")


# Wrapper for schedule to support async functions
def job():
    asyncio.run(send_message())


# Schedule the task to run every 5 seconds
schedule.every(5).seconds.do(job)

# Main loop to run the scheduler
if __name__ == "__main__":
    print("Starting message client...")
    while True:
        schedule.run_pending()
        time.sleep(1)
