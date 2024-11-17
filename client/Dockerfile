# Base image
FROM python:3.9-slim

# Set working directory
WORKDIR /app

# Install dependencies
RUN pip install aiohttp schedule

# Copy application code
COPY app.py /app/

EXPOSE 8080

# Command to run the client
CMD ["python", "app.py"]
