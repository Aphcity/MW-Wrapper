import copy
import os
import time
import re
from contextlib import closing
import threading
import requests
from pydub import AudioSegment
from lxml import etree
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36'
}
failsound = set()
# 输出文件夹
out_dir = "sound_new"
test_dir = "sound_had"
# 线程数

# http请求超时设置
timeout = 200
new = 0


def download(soundName, url):
    global new, allsum
    outpath = os.path.join(out_dir, soundName)
    testpath = os.path.join(test_dir, soundName)
    allsum += 1
    if os.path.isfile(testpath):
        return
    if os.path.isfile(outpath):
        return
    new += 1
    # print("true")
    print("new", new, soundName, url)

    with closing(requests.get(url, stream=True, headers=headers, timeout=timeout)) as r:
        rc = r.status_code
        if 299 < rc or rc < 200:
            print('returnCode%s\t%s' % (rc, url))
            failsound.add(soundName)
            return
        content_length = int(r.headers.get('content-length', '0'))
        if content_length == 0:
            print('size0\t%s' % url)
            return
        with open(outpath, 'wb') as f:
            for data in r.iter_content(1024):
                f.write(data)


def generator():
    global allsum, new
    fpath = open(r'soundUrl.txt', "r")
    text_line = fpath.readline()
    while text_line != '':
        text_line = text_line.strip('\n')
        aa = text_line.split('.mp3_')[0]
        soundName = text_line[:len(aa)+4]
        url = text_line[len(soundName)+1:]
        yield soundName, url
        text_line = fpath.readline()
    fpath.close()


lock = threading.Lock()


def loop(sounds):
    global new, allsum
    allsum = 0
    new = 0

    while True:
        try:
            with lock:
                soundName, url = next(sounds)
        except StopIteration:
            break
        download(soundName, url)


sounds = generator()
failsound = set()
thread_num = 200
threadList = []
for i in range(0, thread_num):
    xx = threading.Thread(target=loop, name='LoopThread %s' %
                          i, args=(sounds,))
    threadList.append(xx)
for y in threadList:
    y.start()
for y in threadList:
    y.join()
failoedf = open('soundfailoed.txt', 'w')
for xx in failsound:
    failoedf.write(xx+"\n")
failoedf.close()

# soundSet =set()
# fpath = open(r'soundUrl.txt',"r")
# text_line = fpath.readline()
# while text_line != '':
#     text_line = text_line.strip('\n')
#     aa = text_line.split('.mp3_')[0]
#     soundName = text_line[:len(aa)+4]
#     soundSet.add(soundName)
#     text_line = fpath.readline()
# fpath.close()
# aa=0
# for xx in os.listdir(out_dir):
#     if xx not in soundSet:
#         outpath = os.path.join(out_dir, xx)
#         os.remove(outpath)
#         aa+=1
#         print(outpath,"removed",aa)
