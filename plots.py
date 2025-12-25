import re
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np

# Установка стиля для красивых графиков
plt.style.use('seaborn-v0_8')
sns.set_palette("husl")

# Чтение файла
with open("bench_vclock.txt", "r", encoding="utf-8", errors="ignore") as f:
    text = f.read().splitlines()

rows = []
for line in text:
    if not line.startswith("Benchmark"):
        continue

    # BenchmarkName-8   iters   ns/op   [optional ...]
    m = re.match(r"^(Benchmark[^\s]+)\s+(\d+)\s+([\d\.]+)\s+ns/op(.*)$", line)
    if not m:
        continue

    name = m.group(1)
    iters = int(m.group(2))
    nsop = float(m.group(3))
    tail = m.group(4)

    # Извлечение B/op и allocs/op
    b_op = None
    allocs_op = None
    m_mem = re.search(r"(\d+)\s+B/op\s+(\d+)\s+allocs/op", tail)
    if m_mem:
        b_op = int(m_mem.group(1))
        allocs_op = int(m_mem.group(2))

    # Извлечение avg_latency_us
    avg_latency_us = None
    m_lat = re.search(r"([\d\.]+)\s+avg_latency_us", tail)
    if m_lat:
        avg_latency_us = float(m_lat.group(1))

    rows.append({
        "name": name,
        "iters": iters,
        "ns_op": nsop,
        "op_s": 1e9 / nsop,
        "b_op": b_op,
        "allocs_op": allocs_op,
        "avg_latency_us": avg_latency_us,
    })

df = pd.DataFrame(rows)

# Группировка и усреднение по name
df_agg = df.groupby("name").agg({
    "op_s": "mean",
    "b_op": "mean",
    "allocs_op": "mean",
    "avg_latency_us": "mean",
    "ns_op": "mean",
}).reset_index()

# Категоризация бенчмарков
def categorize(name):
    if "SingleWrite" in name:
        return "Single Write"
    elif "SingleRead" in name:
        return "Single Read"
    elif "ReadWriteMix" in name:
        return "Read-Write Mix"
    elif "ConcurrentWrites" in name:
        return "Concurrent Writes"
    elif "ConcurrentReads" in name:
        return "Concurrent Reads"
    elif "Transaction" in name:
        return "Transaction"
    elif "HotKey" in name:
        return "Hot Key"
    elif "QuorumOverhead" in name:
        return "Quorum Overhead"
    elif "VectorClockSize" in name:
        return "Vector Clock Size"
    elif "LatencyAvg" in name:
        return "Latency Avg"
    else:
        return "Other"

df_agg["category"] = df_agg["name"].apply(categorize)

# Функция для создания красивых барплотов
def plot_bar(data, x, y, title, xlabel, ylabel, filename, sort_by=None, ascending=False, top_n=20):
    if sort_by:
        data = data.sort_values(sort_by, ascending=ascending).head(top_n)
    fig, ax = plt.subplots(figsize=(12, 8))
    bars = ax.barh(data[x], data[y], color=sns.color_palette("husl", len(data)))
    ax.set_title(title, fontsize=16, fontweight='bold')
    ax.set_xlabel(xlabel, fontsize=12)
    ax.set_ylabel(ylabel, fontsize=12)
    ax.invert_yaxis()
    plt.tight_layout()
    plt.savefig(filename, dpi=300, bbox_inches='tight')
    plt.close()

# 1. Throughput (ops/sec)
plot_bar(df_agg, "name", "op_s", "Throughput Comparison", "Operations per Second", "Benchmark", "throughput.png", sort_by="op_s", ascending=False)

# 2. Memory Usage (B/op)
df_mem = df_agg[df_agg["b_op"].notna()]
if not df_mem.empty:
    plot_bar(df_mem, "name", "b_op", "Memory Usage per Operation", "Bytes per Operation", "Benchmark", "memory_b_op.png", sort_by="b_op", ascending=False)

# 3. Allocations per Operation
df_alloc = df_agg[df_agg["allocs_op"].notna()]
if not df_alloc.empty:
    plot_bar(df_alloc, "name", "allocs_op", "Allocations per Operation", "Allocations per Operation", "Benchmark", "allocs_op.png", sort_by="allocs_op", ascending=False)

# 4. Latency (если есть)
df_lat = df_agg[df_agg["avg_latency_us"].notna()]
if not df_lat.empty:
    plot_bar(df_lat, "name", "avg_latency_us", "Average Latency", "Average Latency (μs)", "Benchmark", "latency.png", sort_by="avg_latency_us", ascending=True)

# 5. Групповой график по категориям
categories = df_agg["category"].unique()
fig, axes = plt.subplots(len(categories), 1, figsize=(12, 6*len(categories)))
if len(categories) == 1:
    axes = [axes]

for i, cat in enumerate(categories):
    sub_df = df_agg[df_agg["category"] == cat].sort_values("op_s", ascending=False)
    if not sub_df.empty:
        axes[i].barh(sub_df["name"], sub_df["op_s"], color=sns.color_palette("husl", len(sub_df)))
        axes[i].set_title(f"{cat} Benchmarks - Throughput", fontsize=14, fontweight='bold')
        axes[i].set_xlabel("Operations per Second")
        axes[i].invert_yaxis()

plt.tight_layout()
plt.savefig("categories_throughput.png", dpi=300, bbox_inches='tight')
plt.close()

# Экспорт в CSV
df_agg.to_csv("bench_vclock_parsed.csv", index=False)

print("Saved: throughput.png, memory_b_op.png, allocs_op.png, latency.png, categories_throughput.png, bench_vclock_parsed.csv")