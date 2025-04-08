# Uriel - AI-Powered API Security Scanner

**Uriel** is a fast and developer-friendly CLI tool that leverages **Llama 3** to scan and analyze API endpoints for potential security flaws. It combines traditional API testing with the power of AI to give you actionable insights and prevent vulnerabilities before they become a threat.

>  Built in Go — Powered by AI

---

## Features

-  **Scan any REST API endpoint**
-  **Llama 3-powered vulnerability analysis** using [Ollama](https://ollama.com/)
-  Detects common security issues: 
  - Missing auth headers
  - Insecure HTTP methods
  - Open endpoints
  - Rate limiting & more
- AI-generated suggestions to fix the flaws
- Clean CLI experience with `cobra`

---

##  Installation

### From source

```bash
git clone https://github.com/your-username/uriel.git
cd uriel
go build -o uriel