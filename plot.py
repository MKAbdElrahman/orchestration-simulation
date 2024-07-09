import pandas as pd
import matplotlib.pyplot as plt

# Load data from CSV files
file_path = 'sweep_results_non_blocking.csv'
file_path_2 = 'sweep_results_blocking.csv'

data = pd.read_csv(file_path)
data_2 = pd.read_csv(file_path_2)

# Convert latency columns to float values representing milliseconds for both datasets
data['AverageSuccessLatency'] = data['AverageSuccessLatency'].str.replace('ms', '').astype('float')
data['AverageFailureLatency'] = data['AverageFailureLatency'].str.replace('ms', '').astype('float')
data_2['AverageSuccessLatency'] = data_2['AverageSuccessLatency'].str.replace('ms', '').astype('float')
data_2['AverageFailureLatency'] = data_2['AverageFailureLatency'].str.replace('ms', '').astype('float')

# Convert percentage strings to float values for both datasets
data['SuccessRatio'] = data['SuccessRatio'].str.rstrip('%').astype('float') / 100.0
data['FailureRatio'] = 1.0 - data['SuccessRatio']
data_2['SuccessRatio'] = data_2['SuccessRatio'].str.rstrip('%').astype('float') / 100.0
data_2['FailureRatio'] = 1.0 - data_2['SuccessRatio']

# Plot combined latencies for both datasets in non-blocking mode
plt.ion()  # Turn on interactive mode (non-blocking)
plt.figure(figsize=(10, 6))
plt.plot(data['Availability'], data['AverageSuccessLatency'], marker='o', linestyle='-', color='c', label='Average Success Latency (Non-blocking)')
plt.plot(data['Availability'], data['AverageFailureLatency'], marker='o', linestyle='-', color='m', label='Average Failure Latency (Non-blocking)')
plt.plot(data_2['Availability'], data_2['AverageSuccessLatency'], marker='x', linestyle='--', color='b', label='Average Success Latency (Blocking)')
plt.plot(data_2['Availability'], data_2['AverageFailureLatency'], marker='x', linestyle='--', color='r', label='Average Failure Latency (Blocking)')
plt.title('Latencies per Availability')
plt.xlabel('Availability')
plt.ylabel('Latency [ms]')
plt.legend()
plt.grid(True)
# plt.savefig('combined_latency_plot.png')
plt.show()

# Plot combined Failure and Success Ratios for both datasets in blocking mode
plt.ioff()  # Turn off interactive mode (blocking)
plt.figure(figsize=(10, 6))
plt.plot(data['Availability'], data['SuccessRatio'], marker='o', linestyle='-', color='y', label='Success Ratio (Non-blocking)')
plt.plot(data_2['Availability'], data_2['SuccessRatio'], marker='x', linestyle='--', color='g', label='Success Ratio (Blocking)')
plt.title('Success Ratios per Availability')
plt.xlabel('Availability')
plt.ylabel('Ratio')
plt.legend()
plt.grid(True)
# plt.savefig('failure_success_ratio_plot.png')
plt.show()  # This will block the execution until the plot window is closed
