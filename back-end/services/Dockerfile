# Use a lightweight Ubuntu base image
FROM ubuntu:22.04

# Install curl and other dependencies
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Install Ollama
RUN curl -fsSL https://ollama.com/install.sh | sh

# Set Ollama to listen on all interfaces, not just localhost
ENV OLLAMA_HOST=0.0.0.0:11434

# Expose the port Ollama uses
EXPOSE 11434

# Run Ollama serve and pull the nomic-embed-text model
CMD ["/bin/bash", "-c", "ollama serve & sleep 10 && ollama pull nomic-embed-text && tail -f /dev/null"]
