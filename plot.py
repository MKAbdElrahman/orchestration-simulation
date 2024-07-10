import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# List of filenames to load
filenames = [
    'epic_3_services.csv',
    'fantasyfiction_3_services.csv',
    'fairytale_3_services.csv',
    'parallel_3_services.csv',

]

# Function to preprocess the data
def preprocess_data(data):
    data['AverageSuccessLatency'] = data['AverageSuccessLatency'].str.replace('ms', '').astype('float')
    data['AverageFailureLatency'] = data['AverageFailureLatency'].str.replace('ms', '').astype('float')
    data['SuccessRatio'] = data['SuccessRatio'].str.rstrip('%').astype('float') / 100.0
    data['FailureRatio'] = 1.0 - data['SuccessRatio']
    return data

# Function to extract parameters from filename
def extract_params(filename):
    parts = filename.split('_')
    type_of_sweep = parts[0]
    num_services = int(parts[-2])
    return type_of_sweep, num_services

# Load and preprocess data from each file
data_frames = {}
for filename in filenames:
    data = pd.read_csv(filename)
    data = preprocess_data(data)
    data_frames[filename] = data

# Get unique number of domain services and colors
service_groups = {}
for filename in filenames:
    type_of_sweep, num_services = extract_params(filename)
    if num_services not in service_groups:
        service_groups[num_services] = []
    service_groups[num_services].append((filename, type_of_sweep))

# Set a vibrant color palette and markers
markers = ['o', 'x']

# Set the aesthetic style of the plots
sns.set(style="whitegrid", font_scale=1.1)

# Function to plot latencies
def plot_latencies(data_frames, title):
    plt.figure(figsize=(10, 6))
    for num_services, files in service_groups.items():
        for filename, type_of_sweep in files:
            data = data_frames[filename]
            linestyle = '-' if 'nonblocking' in type_of_sweep else '--'
            plt.plot(data['Availability'], data['AverageSuccessLatency'], linestyle=linestyle, marker=markers[0], markersize=5, label=f'Success Latency ({type_of_sweep}, {num_services} services)')
            plt.plot(data['Availability'], data['AverageFailureLatency'], linestyle=linestyle, marker=markers[1], markersize=5, label=f'Failure Latency ({type_of_sweep}, {num_services} services)', alpha=0.6)

    plt.title(title, fontsize=16, weight='bold')
    plt.xlabel('Availability', fontsize=12)
    plt.ylabel('Latency [ms]', fontsize=12)
    plt.legend(loc='best', fontsize=6)  # Smaller font size for legend
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.tight_layout()
    plt.show(block=False)  # Show plot non-blocking

# Function to plot success ratios
def plot_success_ratios(data_frames, title):
    plt.figure(figsize=(10, 6))
    color_index = 0
    for num_services, files in service_groups.items():
        for filename, type_of_sweep in files:
            data = data_frames[filename]
            linestyle = '-' if 'nonblocking' in type_of_sweep else '--'
            plt.plot(data['Availability'], data['SuccessRatio'], linestyle=linestyle, marker=markers[0], markersize=5, label=f'Success Ratio ({type_of_sweep}, {num_services} services)')
        color_index += 1

    plt.title(title, fontsize=16, weight='bold')
    plt.xlabel('Availability', fontsize=12)
    plt.ylabel('Success Ratio', fontsize=12)
    plt.legend(loc='best', fontsize=6)  # Smaller font size for legend
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.tight_layout()
    plt.show()

# Plot latencies
plot_latencies(data_frames, 'Latencies per Availability')

# Plot success ratios
plot_success_ratios(data_frames, 'Success Ratios per Availability')

# Keep the plots displayed until explicitly closed
plt.show()
