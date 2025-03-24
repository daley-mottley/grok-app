package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "github.com/joho/godotenv" // For loading .env file
)

// GrokClient handles interactions with the Grok API
type GrokClient struct {
    apiKey  string
    baseURL string
    client  *http.Client
}

// NewGrokClient initializes a new Grok API client
func NewGrokClient() *GrokClient {
    // Load .env file if it exists
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
        // Continue; we’ll fall back to os.Getenv
    }

    apiKey := os.Getenv("GROK_API_KEY")
    if apiKey == "" {
        log.Fatal("GROK_API_KEY environment variable is not set")
    }
    return &GrokClient{
        apiKey:  apiKey,
        baseURL: "https://api.grok.ai/v1", // Hypothetical URL
        client:  &http.Client{},
    }
}

// makeRequest sends an HTTP request to the Grok API
func (c *GrokClient) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
    url := c.baseURL + endpoint
    var req *http.Request
    var err error

    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            log.Printf("Error marshaling request body: %v", err)
            return nil, fmt.Errorf("internal error preparing request")
        }
        req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
    } else {
        req, err = http.NewRequest(method, url, nil)
    }

    if err != nil {
        log.Printf("Error creating request: %v", err)
        return nil, fmt.Errorf("failed to create request")
    }

    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        log.Printf("Error sending request to %s: %v", url, err)
        return nil, fmt.Errorf("network error: unable to reach API")
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("API returned non-200 status: %d", resp.StatusCode)
        return nil, fmt.Errorf("API error: received status %d", resp.StatusCode)
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return nil, fmt.Errorf("failed to read API response")
    }
    return data, nil
}

// AskQuestion sends a question to the Grok API
func (c *GrokClient) AskQuestion(question string) (string, error) {
    body := map[string]string{"question": question}
    data, err := c.makeRequest("POST", "/ask", body)
    if err != nil {
        return "", err
    }
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        log.Printf("Error parsing ask response: %v", err)
        return "", fmt.Errorf("invalid API response")
    }
    answer, ok := result["answer"].(string)
    if !ok {
        log.Println("Answer field not found in response")
        return "", fmt.Errorf("no answer provided by API")
    }
    return answer, nil
}

// GenerateImage sends a prompt for image generation
func (c *GrokClient) GenerateImage(prompt string) (string, error) {
    body := map[string]string{"prompt": prompt}
    data, err := c.makeRequest("POST", "/image", body)
    if err != nil {
        return "", err
    }
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        log.Printf("Error parsing image response: %v", err)
        return "", fmt.Errorf("invalid API response")
    }
    url, ok := result["image_url"].(string)
    if !ok {
        log.Println("Image URL not found in response")
        return "", fmt.Errorf("no image URL provided by API")
    }
    return url, nil
}

// AnalyzeData performs basic analysis on a comma-separated list of numbers
func (c *GrokClient) AnalyzeData(data string) (map[string]float64, error) {
    body := map[string]string{"data": data}
    respData, err := c.makeRequest("POST", "/analyze", body)
    if err != nil {
        return nil, err
    }
    var result map[string]interface{}
    if err := json.Unmarshal(respData, &result); err != nil {
        log.Printf("Error parsing analyze response: %v", err)
        return nil, fmt.Errorf("invalid API response")
    }
    analysis := make(map[string]float64)
    for key, value := range result {
        if num, ok := value.(float64); ok {
            analysis[key] = num
        } else {
            log.Printf("Invalid analysis value for %s: %v", key, value)
            return nil, fmt.Errorf("invalid analysis data")
        }
    }
    return analysis, nil
}

func main() {
    // Set up logging to a file
    logFile, err := os.OpenFile("grok_app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
        os.Exit(1)
    }
    defer logFile.Close()
    log.SetOutput(logFile)
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Initialize Grok API client
    client := NewGrokClient()

    // Welcome message
    fmt.Println("Welcome to the Grok AI Chatbot!")
    fmt.Println("Commands:")
    fmt.Println("  ask: <question>          - Ask a question")
    fmt.Println("  image: <description>     - Generate an image URL")
    fmt.Println("  analyze: <numbers>       - Analyze comma-separated numbers (e.g., 1,2,3)")
    fmt.Println("  exit                     - Quit the app")

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        input := strings.TrimSpace(scanner.Text())
        if input == "exit" {
            fmt.Println("Goodbye!")
            break
        }

        parts := strings.SplitN(input, ": ", 2)
        if len(parts) != 2 && parts[0] != "exit" {
            fmt.Println("Invalid format. Use <command>: <data> or 'exit'")
            continue
        }

        if len(parts) == 1 {
            continue
        }

        command, data := parts[0], parts[1]
        switch strings.ToLower(command) {
        case "ask":
            answer, err := client.AskQuestion(data)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                continue
            }
            fmt.Printf("Answer: %s\n", answer)

        case "image":
            url, err := client.GenerateImage(data)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                continue
            }
            fmt.Printf("Image URL: %s\n", url)

        case "analyze":
            // Validate input
            nums := strings.Split(data, ",")
            for _, num := range nums {
                if _, err := strconv.ParseFloat(strings.TrimSpace(num), 64); err != nil {
                    fmt.Println("Error: Invalid number format in data")
                    log.Printf("Invalid number in analyze input: %s", num)
                    continue
                }
            }
            analysis, err := client.AnalyzeData(data)
            if err != nil {
                fmt.Printf("Error: %v\n", err)
                continue
            }
            fmt.Println("Analysis Results:")
            for key, value := range analysis {
                fmt.Printf("  %s: %.2f\n", key, value)
            }

        default:
            fmt.Println("Unknown command. Available: ask, image, analyze, exit")
        }
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Error reading input: %v", err)
        fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
        os.Exit(1)
    }
}
