import time
import os
import subprocess
import datetime

from subprocess import Popen
from sys import platform

if platform == "linux" or platform == "linux2":
    # linux
    pathExe = os.path.join(os.path.dirname(__file__), "./eyeProtect")
elif platform == "darwin":
    # OS X
    print(f"{platform} is not supported")
    exit(1)
elif platform == "win32":
    # Windows...
    pathExe = os.path.join(os.path.dirname(__file__), "./eyeProtect.exe")
else:
    print(f"{platform} is not supported")
    exit(1)

print(f"{platform}: {pathExe}")
print(f"end  : {datetime.datetime.now().time()}")
while True:
    time.sleep(60 * 10)
    p = Popen([pathExe])
    print(f"start: {datetime.datetime.now().time()}")
    try:
        p.wait(30)
    except subprocess.TimeoutExpired:
        p.kill()
    print(f"end  : {datetime.datetime.now().time()}")
