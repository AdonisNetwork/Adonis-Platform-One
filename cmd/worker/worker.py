import os
import time
import json
import psycopg2
import redis
import requests
import uuid

DB_DSN = os.getenv("DB_DSN")
REDIS_ADDR = os.getenv("REDIS_ADDR", "redis:6379")
GROQ_API_KEY = os.getenv("GROQ_API_KEY")
GROQ_MODEL = os.getenv("GROQ_MODEL", "mixtral-8x7b-32768")  # مثال

r = redis.Redis.from_url(f"redis://{REDIS_ADDR}")

def get_db_conn():
    return psycopg2.connect(DB_DSN)

def call_groq(prompt: str) -> str:
    """
    ساده: یک call به Groq LLM – بسته به API نهایی، این را تنظیم می‌کنی.
    """
    url = "https://api.groq.com/openai/v1/chat/completions"
    headers = {
        "Authorization": f"Bearer {GROQ_API_KEY}",
        "Content-Type": "application/json",
    }
    payload = {
        "model": GROQ_MODEL,
        "messages": [
            {"role": "system", "content": "You are a precise research assistant."},
            {"role": "user", "content": prompt},
        ],
        "temperature": 0.2,
    }
    resp = requests.post(url, headers=headers, json=payload, timeout=60)
    resp.raise_for_status()
    data = resp.json()
    return data["choices[0][message][content]"] if False else data["choices"][0]["message"]["content"]

def run_research_agent(input_text: str):
    """
    Multi-step pipeline: Understanding → Plan → Research → Structure (JSON + Markdown)
    """

    # Step 1: Understanding & Reframing
    understanding_prompt = f"""
You are a research planner. The user request is:

\"\"\"{input_text}\"\"\" 

1) Rewrite it as a clear research question.
2) Identify 3–5 key sub-questions.
3) Identify best sources (regulators, scientific, industry).
Respond in JSON with fields: research_question, sub_questions, sources.
"""
    understanding_raw = call_groq(understanding_prompt)

    try:
        plan = json.loads(understanding_raw)
    except Exception:
        # fallback: wrap raw in JSON
        plan = {
            "research_question": input_text,
            "sub_questions": [],
            "sources": [],
            "raw": understanding_raw,
        }

    # Step 2: Main Research & Synthesis (still single LLM call in MVP)
    synthesis_prompt = f"""
You are a senior research analyst.

Research question:
{plan.get("research_question", input_text)}

Sub-questions:
{json.dumps(plan.get("sub_questions", []), ensure_ascii=False)}

Sources (preferred types or organizations):
{json.dumps(plan.get("sources", []), ensure_ascii=False)}

Write a structured research report with these sections:
1. Overview
2. Key Findings
3. Entities / Devices / Stakeholders (if any)
4. Regulatory or Strategic Trends
5. Risks & Limitations
6. Actionable Insights for decision-makers

Return the answer in TWO parts:

[MARKDOWN]
<full markdown report here>

[JSON]
A minified JSON object with:
- title
- sections (id, title, content)
- insights (list of strings)
- risks (list of strings)
- citations (list of objects: source, url, snippet if known)
"""
    synthesis_raw = call_groq(synthesis_prompt)

    # ساده: split بین [MARKDOWN] و [JSON]
    markdown_part = ""
    json_part = {}

    if "[JSON]" in synthesis_raw:
        md_part, js_part = synthesis_raw.split("[JSON]", 1)
        markdown_part = md_part.replace("[MARKDOWN]", "").strip()
        try:
            json_part = json.loads(js_part.strip())
        except Exception:
            json_part = {"raw": js_part.strip()}
    else:
        markdown_part = synthesis_raw
        json_part = {"raw": synthesis_raw}

    return markdown_part, json_part

def process_job(job_id: str):
    conn = get_db_conn()
    conn.autocommit = True
    cur = conn.cursor()

    # mark as running & get input
    cur.execute("UPDATE jobs SET status='running' WHERE id=%s", (job_id,))
    cur.execute("SELECT input_text FROM jobs WHERE id=%s", (job_id,))
    row = cur.fetchone()
    if not row:
        conn.close()
        return
    input_text = row[0]

    try:
        markdown, json_result = run_research_agent(input_text)
        cur.execute(
            "UPDATE jobs SET status='completed', result_markdown=%s, result_json=%s, error_text=NULL WHERE id=%s",
            (markdown, json.dumps(json_result, ensure_ascii=False), job_id),
        )
    except Exception as e:
        cur.execute(
            "UPDATE jobs SET status='failed', error_text=%s WHERE id=%s",
            (str(e), job_id),
        )
    finally:
        conn.close()

def main_loop():
    print("Worker started, listening for jobs...")
    while True:
        # BRPOP → queue name, timeout=0 → block
        _, job_id_bytes = r.brpop("jobs_queue")
        job_id = job_id_bytes.decode("utf-8")
        print("Processing job:", job_id)
        process_job(job_id)

if __name__ == "__main__":
    main_loop()
