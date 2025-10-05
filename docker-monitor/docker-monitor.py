import time
import docker
import pandas as pd
import matplotlib.pyplot as plt
import sys

client = docker.from_env()
data = []
INTERVAL = 10

if len(sys.argv) < 2:
    print("Uso: python docker_monitor.py <nome_do_container>")
    sys.exit(1)

container_name = sys.argv[1]

try:
    while True:
        try:
            container = client.containers.get(container_name)
            stats = container.stats(stream=False)

            cpu_delta = stats['cpu_stats']['cpu_usage']['total_usage'] - stats['precpu_stats']['cpu_usage']['total_usage']
            system_cpu_delta = stats['cpu_stats']['system_cpu_usage'] - stats['precpu_stats']['system_cpu_usage']
            num_cpus = stats['cpu_stats']['online_cpus']
            cpu_usage = (cpu_delta / system_cpu_delta) * num_cpus * 100 if system_cpu_delta > 0 else 0

            mem_usage = stats['memory_stats']['usage']
            mem_limit = stats['memory_stats']['limit']
            mem_percentage = (mem_usage / mem_limit) * 100 if mem_limit > 0 else 0

            print(f"CPU: {cpu_usage}")
            print(f"Memória: {mem_percentage}")

            data.append({
                'timestamp': time.time(),
                'cpu_usage': cpu_usage,
                'mem_percentage': mem_percentage
            })

        except docker.errors.NotFound:
            print(f"Contêiner '{container_name}' não encontrado.")
            break

        time.sleep(INTERVAL)

except KeyboardInterrupt:
    print("Coleta interrompida.")

df = pd.DataFrame(data)

df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')

fig, ax1 = plt.subplots(figsize=(12, 6))

# Eixo y para CPU
ax1.set_xlabel('Tempo')
ax1.set_ylabel('Uso de CPU (%)', color='blue')
ax1.plot(df['timestamp'], df['cpu_usage'], label='Uso de CPU (%)', color='blue', linestyle='-')
ax1.tick_params(axis='y', labelcolor='blue')

# Cria um segundo eixo y para a memória
ax2 = ax1.twinx()
ax2.set_ylabel('Uso de Memória (%)', color='orange')
ax2.plot(df['timestamp'], df['mem_percentage'], label='Uso de Memória (%)', color='orange', linestyle='--')
ax2.tick_params(axis='y', labelcolor='orange')

plt.title(f'Uso de CPU e memória: {container_name}')
ax1.legend(loc='upper left')
ax2.legend(loc='upper right')
plt.grid()
plt.xticks(rotation=45)
plt.tight_layout()

plt.savefig(f'{container_name}_cpu_mem_usage.png')