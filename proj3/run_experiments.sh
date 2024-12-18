#!/bin/bash
#
#SBATCH --mail-user=jeremyyawei@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj3 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/jeremyyawei/Parallel_Programming/project-3-Jeremytsai6987/proj3/
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=60:00


# Output directory
output_dir="results"
mkdir -p $output_dir

python3 generate_time_plots.py

if [ $? -eq 0 ]; then
    echo "Python script executed successfully. Check the $output_dir directory for results and plots."
else
    echo "Python script execution failed. Please check for errors."
fi
