import subprocess
import time
import csv
import os
import matplotlib.pyplot as plt

# 定義參數
portfolio_file = "data/portfolio_test.csv"
strike_price = 700
risk_free_rate = 0.05
time_to_maturity = 1
steps = 252
threads = [2, 4, 6, 8, 12]
simulations_list = [100000, 1000000, 10000000]
output_dir = "results"
os.makedirs(output_dir, exist_ok=True)

# CSV 文件保存執行時間
times_csv = os.path.join(output_dir, "execution_times.csv")
with open(times_csv, "w", newline="") as f:
    writer = csv.writer(f)
    writer.writerow(["Simulations", "Mode", "Threads", "Time(ms)"])

# 定義執行和計時函數
def run_and_time(command):
    """執行命令並計算執行時間（以毫秒為單位）"""
    start_time = time.time()
    process = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    end_time = time.time()
    duration = (end_time - start_time) * 1000  # 轉換為毫秒
    return process.returncode, duration

# 執行實驗
results = {}
for simulations in simulations_list:
    print(f"Running experiments for {simulations} simulations...")
    results[simulations] = {}

    # Sequential Execution
    command = f"go run main.go --mode=sequential --portfolio={portfolio_file} --K={strike_price} --r={risk_free_rate} --T={time_to_maturity} --steps={steps} --simulations={simulations}"
    exit_code, duration = run_and_time(command)
    if exit_code == 0:
        results[simulations]["Sequential"] = duration
        with open(times_csv, "a", newline="") as f:
            csv.writer(f).writerow([simulations, "Sequential", 1, duration])
        print(f"Sequential execution completed in {duration:.2f} ms.")
    else:
        print(f"Sequential execution failed for {simulations} simulations.")

    # Parallel Execution
    for t in threads:
        command = f"go run main.go --mode=parallel --portfolio={portfolio_file} --K={strike_price} --r={risk_free_rate} --T={time_to_maturity} --steps={steps} --simulations={simulations} --threads={t}"
        exit_code, duration = run_and_time(command)
        if exit_code == 0:
            results[simulations][f"Parallel-{t}"] = duration
            with open(times_csv, "a", newline="") as f:
                csv.writer(f).writerow([simulations, "Parallel", t, duration])
            print(f"Parallel execution with {t} threads completed in {duration:.2f} ms.")
        else:
            print(f"Parallel execution failed for {t} threads and {simulations} simulations.")

    # Parallel Work-Stealing Execution
    for t in threads:
        command = f"go run main.go --mode=parallel-stealing --portfolio={portfolio_file} --K={strike_price} --r={risk_free_rate} --T={time_to_maturity} --steps={steps} --simulations={simulations} --threads={t}"
        exit_code, duration = run_and_time(command)
        if exit_code == 0:
            results[simulations][f"Parallel-Stealing-{t}"] = duration
            with open(times_csv, "a", newline="") as f:
                csv.writer(f).writerow([simulations, "Parallel-Stealing", t, duration])
            print(f"Parallel Work-Stealing execution with {t} threads completed in {duration:.2f} ms.")
        else:
            print(f"Parallel Work-Stealing execution failed for {t} threads and {simulations} simulations.")

# 繪製加速比圖
for simulations, data in results.items():
    sequential_time = data.get("Sequential", None)
    if not sequential_time:
        print(f"Sequential time missing for {simulations} simulations. Skipping...")
        continue

    # 繪製 Parallel 模式的加速比
    speedups_parallel = []
    speedups_work_stealing = []
    for t in threads:
        parallel_time = data.get(f"Parallel-{t}", None)
        stealing_time = data.get(f"Parallel-Stealing-{t}", None)

        if parallel_time:
            speedups_parallel.append(sequential_time / parallel_time)
        else:
            speedups_parallel.append(0)

        if stealing_time:
            speedups_work_stealing.append(sequential_time / stealing_time)
        else:
            speedups_work_stealing.append(0)

    plt.figure()
    plt.plot(threads, speedups_parallel, marker="o", label="Parallel")
    plt.plot(threads, speedups_work_stealing, marker="x", label="Work-Stealing")
    plt.title(f"Speedup for {simulations} Simulations")
    plt.xlabel("Number of Threads")
    plt.ylabel("Speedup")
    plt.legend()
    plt.grid()
    plt.savefig(os.path.join(output_dir, f"speedup_{simulations}.png"))
    print(f"Speedup plot saved for {simulations} simulations.")

print("All experiments completed.")
