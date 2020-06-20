import urllib3
import random
import socket
import struct
import time
import threading
import numpy as np

def split(arr, size):
    arrs = []
    while len(arr) > size:
        pice = arr[:size]
        arrs.append(pice)
        arr = arr[size:]
    arrs.append(arr)
    return arrs

def random_ip():
    return socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff)))

class Colors:
    @staticmethod
    def blue(msj): return '\033[94m' + str(msj) + '\033[0m'

    @staticmethod
    def green(msj): return '\033[92m' + str(msj) + '\033[0m'

    @staticmethod
    def red(msj): return '\033[91m' + str(msj) + '\033[0m'

    @staticmethod
    def bold(msj): return '\033[1m' + str(msj) + '\033[0m'

class Client(threading.Thread):
    def __init__(self, client, ip):
        super(Client, self).__init__()
        self.client = client
        self.ip = ip
        self.result = None

    def run(self):
        target = 'http://localhost:8080/user'
        start = time.time()
        response = self.client.request(
            'POST', target,
                headers={'Content-Type': 'application/json'},
                body="{\"ip\": \""+random_ip()+"\"}")
        if response.status == 200:
            time_measure = round(time.time() - start, 2)
            print(self.ip + " --> " + Colors.green("OK ") + "in " + Colors.bold(str(time_measure)) + "s")
            self.result = time_measure
        else:
            print(Colors.red("ERROR") + ": " + response.reason)

class Main:    
    def __init__(self):
        client = urllib3.PoolManager(10, headers={'user-agent': 'Mozilla/5.0 (Windows NT 6.3; rv:36.0) ..'})
        seconds = int(input(Colors.bold('seconds: ')))
        parallelism = int(input(Colors.bold('parallelism: ')))

        time_measurments = []
        failures = 0
        
        start = time.time()
        while time.time() - start < seconds:
            requests = []
            for _ in range(0, 100000):
                requests.append(Client(client, random_ip()))

            for treads in split(requests, parallelism):
                for thread in treads:
                    thread.start()
                for thread in treads:
                    thread.join()
                for thread in treads:
                    if thread.result is not None:
                        time_measurments.append(thread.result)
                    else:
                        failures+=1

        failure_str = Colors.green(str(failures)+ " failures") if failures == 0 else Colors.red(str(failures)+ " failures")
        print("\nPerformed "+ Colors.bold(str(len(time_measurments))) +" requests in: " + Colors.bold(str(round(np.sum(time_measurments), 2))) + "s With "+failure_str)
        print("Time measurement percentiles:")
        print("p50: " + Colors.bold(str(round(np.percentile(time_measurments, 50), 2)))+"s")
        print("p70: " + Colors.bold(str(round(np.percentile(time_measurments, 70), 2)))+"s")
        print("p90: " + Colors.bold(str(round(np.percentile(time_measurments, 90), 2)))+"s")
        print("p99: " + Colors.bold(str(round(np.percentile(time_measurments, 99), 2)))+"s")

    @staticmethod
    def split(arr, size):
        arrs = []
        while len(arr) > size:
            pice = arr[:size]
            arrs.append(pice)
            arr = arr[size:]
        arrs.append(arr)
        return arrs

    @staticmethod
    def random_ip():
        return socket.inet_ntoa(struct.pack('>I', random.randint(1, 0xffffffff)))

if __name__ == '__main__':
    Main()