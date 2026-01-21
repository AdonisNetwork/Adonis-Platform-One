import os
import time
import json
import psycopg2
import traceback
import redis
import uuid

# -------------------------------------------------------------------
#  Worker Configuration
# -------------------------------------------------------------------

DB_DSN = os.getenv("DB_DSN")
REDIS_ADDR = os.getenv("REDIS_ADDR", "redis://localhost:6379")
QUEUE_NAME = os.getenv("JOBS_QUEUE_NAME", "jobs_queue")
POLL_INTERVAL = float(os.getenv("WORKER_POLL_INTERVAL", "0.5"))

# -------------------------------------------------------------------
#  Database & Redis Connections
# -------------------------------------------------------------------

def connect_postgres():
    return psycopg2.connect(DB_DSN)

def connect_redis():
    return redis.Redis.from_url(REDIS_ADDR, decode_responses=True)

# -------------------------------------------------------------------
#  Event Logger
# -------------------------------------------------------------------

def log_event(conn, job_id, event_type, payload=None):
    with conn.cursor() as cur:
        cur.execute(
            """
            INSERT INTO job_events (job_id, event_type, payload)
            VALUES (%s, %s, %s)
            """,
            (str(job_id), event_type, json.dumps(payload) if payload else None)
        )
    conn.commit()

# -------------------------------------------------------------------
#  Job Status Updater
# -------------------------------------------------------------------

def update_status(conn, job_id, status, status_reason=None, result_markdown=None, result_json=None, error_text=None):
    with conn.cursor() as cur:
        cur.execute(
            """
            UPDATE jobs
            SET status = %s,
                status_reason = %s,
                result_markdown = %s,
                result_json = %s,
                error_text = %s,
                updated_at = NOW()
            WHERE id = %s
            """,
            (
                status,
                status_reason,
                result_markdown,
                json.dumps(result_json) if result_json else None,
                error_text,
                str(job_id),
            ),
        )
    conn.commit()

# -------------------------------------------------------------------
#  AI Task Execution (MVP version)
#   Groq / OpenAI / Local LLM / Agents
# -------------------------------------------------------------------

def execute_task(task_type, input_text, input_meta):
    """
    MVP:
        - multi-agent A1 runtime
        - Groq LPU inference
        - OpenAI GPT-4.1
        - Local models (edge)
        - Tool calling
    """

    return {
        "markdown": f"### Task: {task_type}\n\nInput:\n```\n{input_text}\n```",
        "json": {
            "task": task_type,
            "input": input_text,
            "meta": input_meta,
            "processed_at": time.time(),
        }
    }

# -------------------------------------------------------------------
#  Worker Loop
# -------------------------------------------------------------------

def worker_loop():
    redis_client = connect_redis()

    print("[A1 Worker] Connected to Redis:", REDIS_ADDR)
    print("[A1 Worker] Listening on queue:", QUEUE_NAME)

    while True:
        try:
            job_id = redis_client.rpop(QUEUE_NAME)

            if not job_id:
                time.sleep(POLL_INTERVAL)
                continue

            print(f"[A1 Worker] → Received Job: {job_id}")

            conn = connect_postgres()

            # Load job
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT task_type, input_text, input_meta
                    FROM jobs
                    WHERE id = %s
                    """,
                    (job_id,)
                )
                row = cur.fetchone()

            if not row:
                print("[A1 Worker] WARNING: Job not found:", job_id)
                continue

            task_type, input_text, input_meta_raw = row

            input_meta = None
            if input_meta_raw:
                try:
                    input_meta = json.loads(input_meta_raw)
                except:
                    input_meta = {}

            # -------------------------------------------
            # Update → processing
            # -------------------------------------------
            update_status(conn, job_id, "processing", "job started")
            log_event(conn, job_id, "processing", {"info": "Job execution started"})

            # -------------------------------------------
            # Execute Task
            # -------------------------------------------
            try:
                result = execute_task(task_type, input_text, input_meta)

                update_status(
                    conn,
                    job_id,
                    "completed",
                    "job finished",
                    result_markdown=result["markdown"],
                    result_json=result["json"],
                )

                log_event(conn, job_id, "completed", result["json"])

                print(f"[A1 Worker] ✓ Job Completed: {job_id}")

            except Exception as e:
                error_msg = traceback.format_exc()

                update_status(
                    conn,
                    job_id,
                    "failed",
                    "worker error",
                    error_text=error_msg,
                )

                log_event(
                    conn,
                    job_id,
                    "failed",
                    {"error": str(e), "trace": error_msg},
                )

                print(f"[A1 Worker] ✗ Job Failed: {job_id}")
                print(error_msg)

            finally:
                conn.close()

        except Exception as e:
            print("[A1 Worker] CRITICAL ERROR:", e)
            time.sleep(1)

# -------------------------------------------------------------------
#  Main Entry
# -------------------------------------------------------------------

if __name__ == "__main__":
    print("[A1 Worker] Starting...")
    worker_loop()
