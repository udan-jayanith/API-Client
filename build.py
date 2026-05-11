import os
import platform
import subprocess
import sys
from pathlib import Path

file_path = Path(__file__).resolve()
root_dir = file_path.parent


def execute_command(command):
    # Determine the command based on the OS
    if platform.system() == "Windows":
        use_shell = True  # 'dir' is a shell built-in on Windows
    else:
        use_shell = False  # 'ls' is a standalone executable on Unix

    try:
        cmd = subprocess.run(command, capture_output=True, text=True, shell=use_shell)
        return cmd.stdout
    except Exception as e:
        return f"An error occurred: {e}"


def run_zbolt():
    result = execute_command(["go", "run", root_dir])
    if result != "":
        print(result)


def main():
    args = sys.argv
    if len(args) == 2 and args[1] == "run":
        run_zbolt()
        os._exit(0)
    elif len(args) >= 2:
        print("Unknow argument passed")
        os._exit(1)

    print("[1] Genrate\n[2] Build\n[3] Run")
    option = input("Enter: ")
    
    if option == "1":
        run_zbolt()
        os._exit(1)


try:
    main()
except KeyboardInterrupt:
    print("\nScript stopped by Ctrl+C")
