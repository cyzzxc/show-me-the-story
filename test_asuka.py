import requests, time

BASE = "http://localhost:48090"
PROJECT = "test-images"
INTENT = "星穹铁道云璃，动作服装场景自定，@zbjlm 画师画风"
COUNT = 5

def api(method, path, body=None):
    url = f"{BASE}{path}"
    resp = requests.request(method, url, json=body, timeout=300)
    resp.raise_for_status()
    return resp.json()

print(f"intent: {INTENT}")
api("POST", "/api/projects/select", {"name": PROJECT})

# 单次 prepare 批量出 5 个 prompt
t0 = time.time()
r = api("POST", "/api/media/image/prepare", {
    "intent": INTENT, "anima": True, "count": COUNT
})
results = r.get("results", [r])  # count=1 返回单对象，count>1 返回 {"results":[...]}
print(f"prepare done in {time.time()-t0:.0f}s → {len(results)} prompts")
for i, p in enumerate(results):
    print(f"  [{i+1}] {p['prompt'][:60].replace(chr(10),' ')}...")

# 批量 generate
print(f"\n[generate] 批量推送 {len(results)} 个 prompt...")
api("POST", "/api/media/image/generate", {
    "prompts": [{
        "prompt": p["prompt"],
        "negative_prompt": p.get("negative_prompt", ""),
        "resolution": p.get("resolution", "1080*1440"),
        "tags": p.get("tags", []),
        "count": 1,
    } for p in results]
})

while True:
    status = api("GET", "/api/status")
    if not status.get("is_task_running"):
        break
    time.sleep(5)

images = api("GET", "/api/media/images")
cutoff = time.time() - 600
recent = [img for img in images
          if time.mktime(time.strptime(img["created_at"][:19], "%Y-%m-%dT%H:%M:%S")) >= cutoff]
print(f"\n[DONE] 本次生图 {len(recent)} 张:")
for img in recent:
    print(f"  {img['file']}  {img['resolution']}")
