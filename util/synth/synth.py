#!/usr/bin/python3

""" The objective of synth.py is to emulate the behavior of a correctly 
formatted HTTP-post message from a device or formatting function.

Usage:
    synth.py <START> [STEPS] [LOW] [HIGH] [DELAY] [TIMESTAMP] 
"""
from datetime import datetime
import random
import time
import sys

prob = [0.5, 0.5]

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("""Usage:
synth.py <START> [STEPS] [LOW] [HIGH] [DELAY] [TIMESTAMP]""")
        sys.exit(1)

    start = float(sys.argv[1]) 
    low, high = 0, 0
    now = ""
    delay = 1.0

    if len(sys.argv) >= 3:
        steps = int(sys.argv[2])
    else:
        steps = 10

    if len(sys.argv) <= 5:
        low = start * 0.75
        high = start * 1.25
    else:
        low = float(sys.argv[3])
        high = float(sys.argv[4])

    if len(sys.argv) >= 6:
        delay = float(sys.argv[5])

    if len(sys.argv) >= 7:
        now = datetime.now().isoformat()

    # print("dest_addr: %s dest_port: %d start: %s steps: %d low: %d high: %d" % (dest_addr,
    #     dest_port, start, steps, low, high))
    
    for i in range(steps):
        draw = random.uniform(0, 1)
        if draw >= prob[0]:
            start += 1
        else:
            start -= 1
        
        if start >= high: 
            start = high
        if start <= low:
            start = low
        # print("draw:%-6.3f value:%-6.0f" % (draw, start))
        # print("%.0f" % start)
        
        if now != "":
            sys.stdout.write("%s,%.5f\n" % (now, start))
        else:
            sys.stdout.write("%.5f\n" % start)
        sys.stdout.flush()

        time.sleep(delay)


    




        
    