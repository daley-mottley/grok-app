# Grok AI Chatbot

A command-line chatbot built in Go that leverages the hypothetical Grok API from xAI. This app allows users to ask questions, generate image URLs (simulated), and analyze numerical datasets, all with robust error handling and secure API key management.

## Features
- **Ask Questions**: Get answers from the Grok API.
- **Generate Images**: Request image URLs based on textual prompts (simulated responses).
- **Analyze Data**: Perform basic statistical analysis on comma-separated numbers.
- **Robust Error Handling**: Logs errors to a file (`grok_app.log`) and provides user-friendly feedback.
- **Secure API Key Management**: Loads the API key from a `.env` file using `godotenv`.

## Prerequisites
- **Go**: Version 1.21 or later.
- **Git**: For version control (optional).
- **Grok API Key**: Obtainable from [ide.x.ai](https://ide.x.ai/) (see [Obtaining an API Key](#obtaining-an-api-key)).

## Installation

### Step 1: Set Up Your Environment
Ensure Go is installed and configured on your system:
```
go version
```

### Step 2: Clone or Create the Project
Create a project directory and initialize it as a Go module:
```
mkdir grok-app
cd grok-app
go mod init grok-app
```

### Step 3: Install Dependencies
Add the `godotenv` library for `.env` file support:
```
go get github.com/joho/godotenv
```

### Step 4: Add Project Files
- Copy `main.go` from this repository into `grok-app/main.go`.
- Create a `.env` file with your Grok API key:
```
echo "GROK_API_KEY=your_api_key_here" > .env
```
- Add a `.gitignore` file to exclude sensitive files:
```
echo -e "grok-app\ngrok_app.log\n.env\nvendor/\n*.tmp\n*.swp\n.vscode/\n.idea/" > .gitignore
```

### Step 5: Build the App
Compile the app into an executable:
```
go build -o grok-app
```

## Usage
Run the app and use the following commands:
```
./grok-app
```
On Windows, use `grok-app.exe` instead.

### Available Commands
- `ask: <question>` - Ask a question (e.g., `ask: What is the capital of France?`).
- `image: <description>` - Generate an image URL (e.g., `image: A sunset over the ocean`).
- `analyze: <numbers>` - Analyze numbers (e.g., `analyze: 1, 2, 3, 4, 5`).
- `exit` - Quit the app.

### Example Session
```
> ask: What is 2 + 2?
Answer: 4
> image: A fluffy cat
Image URL: https://example.com/fluffy_cat.jpg
> analyze: 1, 2, 3, 4,5
Analysis Results:
  mean: 3.00
  median: 3.00
> exit
Goodbye!
```

## Obtaining an API Key
1. Sign into [ide.x.ai](https://ide.x.ai/) with your X account.
2. Click your username (top-right) and select "API Keys."
3. Click "Create API Key," set permissions (e.g., `chat:write`), and save.
4. Copy the key and store it in your `.env` file as `GROK_API_KEY`.


## Logging
Errors and events are logged to `grok_app.log`. View logs with:
```
cat grok_app.log
```
On Windows, use `type grok_app.log` or a text editor.

## Troubleshooting
- **"GROK_API_KEY not set"**: Ensure `.env` exists and contains your key.
- **Build Errors**: Verify Go is installed (`go version`) and dependencies are fetched (`go mod tidy`).
- **Network Issues**: Check your internet connection.
