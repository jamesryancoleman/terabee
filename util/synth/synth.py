#!/usr/bin/python3

""" The objective of synth.py is to emulate the behavior of a correctly 
formatted HTTP-post message from a device or formatting function.

Usage:
    synth.py <DEST_ADDR> <DEST_PORT> <START> [STEPS] [LOW] [HIGH] [FMT]
"""
from datetime import datetime
import random
import time
import sys

prob = [0.5, 0.5]

def FormatMessage(fmt:str, value:float) -> str:
    now = datetime.now().isoformat()
    if fmt == "mortar":
        fmt_str = """{
    "time": "%s",
    "value": %f,
    "id": 999
}"""
        return fmt_str % (now, value)
    
    elif fmt == "frost":
        fmt_str = """{
    "phenomenonTime": "%s",
    "resultTime": "%s",
    "result": %f,
    "Datastream": {
        "@iot.id": 1
}"""
        return fmt_str % (now, now, value)
    else:
        fmt_str = "%.3f"
        return fmt_str % (value)

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("""Usage:
synth.py <DEST_ADDR> <DEST_PORT> <START> [STEPS] [LOW] [HIGH] [FMT]""")
        sys.exit(1)

    dest_addr = sys.argv[1]
    dest_port = int(sys.argv[2])
    
    start = float(sys.argv[3]) 
    low, high = 0, 0

    fmt = "value"

    if len(sys.argv) >= 5:
        steps = int(sys.argv[4])
    else:
        steps = 10

    if len(sys.argv) >= 7:
        low = start * 0.75
        high = start * 1.25
    else:
        low = float(sys.argv[5])
        high = float(sys.argv[6])

    if len(sys.argv) >= 8:
        fmt = sys.argv[7]

    # print("dest_addr: %s dest_port: %d start: %s steps: %d low: %d high: %d" % (dest_addr,
    #     dest_port, start, steps, low, high))
    
    for i in range(steps):
        draw = random.uniform(0, 1)
        if draw > prob[1]:
            start += 1
        else:
            start -= 1
        
        if start >= high: 
            start = high
        if start <= low:
            start = low
        # print("draw:%-6.3f value:%-6.0f" % (draw, start))
        # print("%.0f" % start)
        print(FormatMessage(fmt, start))
        time.sleep(0.25)

    




        
    