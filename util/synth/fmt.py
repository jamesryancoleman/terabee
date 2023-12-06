#!/usr/bin/python3

"""fmt.py takes a value read from standard input and embeds it in the
format of your choice.

Usage: fmt.py <value|frost|mortar>
"""

from datetime import datetime
import sys

fmt = "value"

def FormatMessage(fmt:str, value:float) -> str:
    now = datetime.now().isoformat()
    if fmt == "mortar":
        fmt_str = """{"time": "%s","value": %f,"id": 999}"""
        return fmt_str % (now, value)
    
    elif fmt == "frost":
        fmt_str = """{"phenomenonTime": "%s","resultTime": "%s","result": %f,"Datastream": {"@iot.id": 1}}"""
        return fmt_str % (now, now, value)
    else:
        fmt_str = "%.3f"
        return fmt_str % (value)
    

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("usage: fmt.py <frost|mortar>")

    fmt = sys.argv[1]

    while True:
        line = sys.stdin.readline()
        if line == "":
            break
        else:
            print(FormatMessage(fmt, float(line.strip())))
    