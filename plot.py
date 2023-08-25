import os
import re
import sys
import numpy as np
import matplotlib.pyplot as plt

data = '''
BenchmarkCounter                              	57502521	        20.83 ns/op
BenchmarkCounter-2                            	33940744	        35.27 ns/op
BenchmarkCounter-4                            	11385031	        98.90 ns/op
BenchmarkCounter-8                            	 8708835	       138.2 ns/op
BenchmarkCounter-16                           	 5410040	       197.2 ns/op
BenchmarkCounter-32                           	 5832330	       214.5 ns/op
BenchmarkCounter-64                           	 5439838	       217.9 ns/op
BenchmarkAtomicCounter                        	137129950	         8.758 ns/op
BenchmarkAtomicCounter-2                      	70709760	        18.68 ns/op
BenchmarkAtomicCounter-4                      	67458786	        18.85 ns/op
BenchmarkAtomicCounter-8                      	46132720	        21.94 ns/op
BenchmarkAtomicCounter-16                     	43168729	        27.46 ns/op
BenchmarkAtomicCounter-32                     	42679546	        27.42 ns/op
BenchmarkAtomicCounter-64                     	42640140	        27.09 ns/op
BenchmarkPartitionedAtomicCounter             	79664224	        14.70 ns/op
BenchmarkPartitionedAtomicCounter-2           	49263518	        25.06 ns/op
BenchmarkPartitionedAtomicCounter-4           	44560155	        26.06 ns/op
BenchmarkPartitionedAtomicCounter-8           	42517705	        35.48 ns/op
BenchmarkPartitionedAtomicCounter-16          	42816031	        28.39 ns/op
BenchmarkPartitionedAtomicCounter-32          	78738529	        13.76 ns/op
BenchmarkPartitionedAtomicCounter-64          	123775238	        11.21 ns/op
BenchmarkPaddedPartitionedAtomicCounter       	79984232	        15.34 ns/op
BenchmarkPaddedPartitionedAtomicCounter-2     	157331673	         7.623 ns/op
BenchmarkPaddedPartitionedAtomicCounter-4     	316654566	         3.799 ns/op
BenchmarkPaddedPartitionedAtomicCounter-8     	590977437	         2.043 ns/op
BenchmarkPaddedPartitionedAtomicCounter-16    	1000000000	         1.011 ns/op
BenchmarkPaddedPartitionedAtomicCounter-32    	1000000000	         0.5740 ns/op
BenchmarkPaddedPartitionedAtomicCounter-64    	1000000000	         0.3851 ns/op
BenchmarkMap                                  	37959387	        30.03 ns/op
BenchmarkMap-2                                	24338346	        49.47 ns/op
BenchmarkMap-4                                	12017052	       107.2 ns/op
BenchmarkMap-8                                	 6072268	       180.9 ns/op
BenchmarkMap-16                               	 5237257	       212.6 ns/op
BenchmarkMap-32                               	 5802340	       203.8 ns/op
BenchmarkMap-64                               	 5671561	       202.1 ns/op
BenchmarkMapRW                                	40482854	        27.59 ns/op
BenchmarkMapRW-2                              	12107698	        96.36 ns/op
BenchmarkMapRW-4                              	11542309	        92.24 ns/op
BenchmarkMapRW-8                              	12642078	       126.9 ns/op
BenchmarkMapRW-16                             	11115829	       103.3 ns/op
BenchmarkMapRW-32                             	12833227	        97.01 ns/op
BenchmarkMapRW-64                             	18793011	        59.56 ns/op
BenchmarkPartitionedMap                       	33812356	        30.14 ns/op
BenchmarkPartitionedMap-2                     	79382323	        15.11 ns/op
BenchmarkPartitionedMap-4                     	72083950	        41.65 ns/op
BenchmarkPartitionedMap-8                     	258519758	         6.492 ns/op
BenchmarkPartitionedMap-16                    	195574123	        10.22 ns/op
BenchmarkPartitionedMap-32                    	114099651	        12.66 ns/op
BenchmarkPartitionedMap-64                    	70103785	        17.47 ns/op
BenchmarkLockFreeMap                          	79386301	        33.52 ns/op
BenchmarkLockFreeMap-2                        	32820781	        33.25 ns/op
BenchmarkLockFreeMap-4                        	80829972	        14.82 ns/op
BenchmarkLockFreeMap-8                        	408663632	         3.622 ns/op
BenchmarkLockFreeMap-16                       	563676435	         3.062 ns/op
BenchmarkLockFreeMap-32                       	1000000000	         1.720 ns/op
BenchmarkLockFreeMap-64                       	1000000000	         0.5590 ns/op
BenchmarkStack                                	49726214	       104.6 ns/op
BenchmarkStack-2                              	15882134	        71.11 ns/op
BenchmarkStack-4                              	 9767264	       127.5 ns/op
BenchmarkStack-8                              	 5977665	       172.5 ns/op
BenchmarkStack-16                             	 4469860	       264.0 ns/op
BenchmarkStack-32                             	 4660720	       222.1 ns/op
BenchmarkStack-64                             	 5233731	       226.5 ns/op
BenchmarkLockFreeStack                        	222342340	       106.8 ns/op
BenchmarkLockFreeStack-2                      	44926143	        27.00 ns/op
BenchmarkLockFreeStack-4                      	12772394	        91.22 ns/op
BenchmarkLockFreeStack-8                      	 7988970	       283.0 ns/op
BenchmarkLockFreeStack-16                     	 3444560	       324.4 ns/op
BenchmarkLockFreeStack-32                     	 2817578	       419.2 ns/op
BenchmarkLockFreeStack-64                     	 2458220	       501.2 ns/op
'''

ds_proc_time: dict[tuple[str, int], float] = {}
for line in data.split('\n'):
    line = line.strip()
    if not line:
        continue
    columns = re.split('\\s+', line)
    name_cpu = columns[0][len('Benchmark'):].split('-')
    name = name_cpu[0]
    cpus = int(name_cpu[1]) if len(name_cpu) > 1 else 1
    t = float(columns[2])
    ds_proc_time[(name, cpus)] = t

def plot(ds_list: list[str]):

    # Data for the grouped bars
    categories = [1, 2, 4, 8, 16, 32, 64]

    # Generate an array of x positions for each group
    x = np.arange(len(categories))

    # Set the width of the bars
    bar_width = 0.9
    ds_bar_width = bar_width/len(ds_list)

    # Create the grouped bar chart
    for i, ds in enumerate(ds_list):
        values = [t for k, t in ds_proc_time.items() if k[0] == ds]
        plt.bar(x + float(ds_bar_width*(i - len(ds_list)/2)), values, width=ds_bar_width, label=ds)

    # Customize the plot
    plt.xlabel("CPUs")
    plt.ylabel("ns / ops")
    plt.title("Performance")
    plt.xticks(x, categories)
    plt.legend()

    plt.tight_layout()

    # Display the plot
    plt.show()

plot(sys.argv[1:])
