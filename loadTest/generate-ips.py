import urllib3
import random
import socket
import struct
import time
import threading
import numpy as np
import gevent
from gevent.pool import Pool

def random_ip():
    return socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff)))

def run(client, rusult_times):        
    start = time.time()
    ip = random_ip()
    response = client.request(
        'POST', 'http://localhost:8080/user',
            headers={'Content-Type': 'application/json'},
            body="{\"ip\": \""+ip+"\"}")

    if response.status == 200 or response.status == 400:
        time_measure = round(time.time() - start, 2)
        print(ip + " --> " + Colors.green("OK ") + "in " + Colors.bold(str(time_measure)) + "s")
        rusult_times.append(time_measure)
    else:
        print(Colors.red("ERROR") + ": " + response.reason)

class Colors:
    @staticmethod
    def blue(msj): return '\033[94m' + str(msj) + '\033[0m'

    @staticmethod
    def green(msj): return '\033[92m' + str(msj) + '\033[0m'

    @staticmethod
    def red(msj): return '\033[91m' + str(msj) + '\033[0m'

    @staticmethod
    def bold(msj): return '\033[1m' + str(msj) + '\033[0m'

class Main:     
    def __init__(self):
        client = urllib3.PoolManager(10, headers={'user-agent': 'Mozilla/5.0 (Windows NT 6.3; rv:36.0) ..'})

        num = int(input(Colors.bold('requests: ')))
        #duration = int(input(Colors.bold('duration (seconds): ')))
        duration = None
        concurrency = int(input(Colors.bold('concurrency: ')))

        failures = 0
        start = time.time()

        pool = Pool(concurrency)
        jobs = None
        time_measurments = []

        try:
            if num is not None:
                jobs = [pool.spawn(run, client, time_measurments) for i in range(num)]
                pool.join()
            else:
                with gevent.Timeout(duration, False):
                    jobs = []
                    while True:
                        jobs.append(pool.spawn(run, client, time_measurments))
                    pool.join()
        except KeyboardInterrupt:
            # In case of a keyboard interrupt, just return whatever already got
            # put into the result object.
            pass
        finally:
            total_time = time.time() - start

        failure_str = Colors.green(str(failures)+ " failures") if failures == 0 else Colors.red(str(failures)+ " failures")
        print("\nPerformed "+ Colors.bold(str(len(time_measurments))) +" requests in: " + Colors.bold(str(round(total_time, 2))) + "s With "+failure_str)
        print("Time measurement percentiles:")
        print("p50: " + Colors.bold(str(round(np.percentile(time_measurments, 50), 2)))+"s")
        print("p70: " + Colors.bold(str(round(np.percentile(time_measurments, 70), 2)))+"s")
        print("p90: " + Colors.bold(str(round(np.percentile(time_measurments, 90), 2)))+"s")
        print("p99: " + Colors.bold(str(round(np.percentile(time_measurments, 99), 2)))+"s")

    @staticmethod
    def random_ip():
        return socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff)))

if __name__ == '__main__':
    Main()