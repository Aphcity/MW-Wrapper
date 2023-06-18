import requests
import os
import json
from contextlib import closing
from lxml import etree


def get_proxy():
    return requests.get("http://127.0.0.1:5010/get/").json()


def delete_proxy(proxy):
    requests.get("http://127.0.0.1:5010/delete/?proxy={}".format(proxy))


def newSession():
    global proxy, rs
    url0 = 'https://www.merriam-webster.com/'
    proxy = get_proxy().get("proxy")
    rs = requests.session()
    rs.headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.42'}
    rs.get(url0, proxies={"http": "http://{}".format(proxy)})
    rs.cookies = requests.utils.cookiejar_from_dict(
        {'mwabtest-2022redesign-phase3': 'cross-dungarees-lite', 'x-stack': 'cross-dungarees-lite'}, cookiejar=None, overwrite=True)


def main():
    global sum, allentry
    sum = 0
    urlHead = 'https://www.merriam-webster.com/browse/thesaurus'
    newSession()
    headWords = [chr(i)
                 for i in range(ord('a'), ord('z')+1)] + ['0', 'bio', 'geo']
    xxset = set()
    allentry = set()
    # newandgot = set()
    # with open('test/haveit.txt', 'r') as fx:
    #     text_line = fx.readline()
    #     while text_line != '':
    #         thename = text_line.strip('\n')
    #         newandgot.add(thename)
    #         text_line = fx.readline()
    # print(len(newandgot))

    for a_head in headWords[0:26]:
        index = 1
        while True:
            url = urlHead + '/' + a_head + '/' + str(index)
            ob = download(url, '#default', t='catalog')
            root = etree.HTML(ob)
            xxs = root.xpath("//a[@class='pb-4 pr-4 d-block']")
            for xx in xxs:
                entry = xx.attrib['href'][11:].replace("/", "_")
                xxhref = 'https://www.merriam-webster.com' + \
                    xx.attrib['href']
                download(xxhref, entry)
            yy = root.xpath("//span[@class='counters']")[0]
            tt = yy.text.split(' ')
            if tt[1] == tt[3]:
                break
            else:
                print(tt[1], '/', tt[3], 'page continue!')
                index += 1


def download(url, entry, t='data'):
    global sum, allentry
    retry = 3
    while retry > 0:
        try:
            if t == 'data':
                outpath = 'E:/Golang/mw/thsr/' + entry + '.html'
                if os.path.isfile(outpath):
                    try:
                        print("had:", entry)
                    except:
                        pass
                    return
                try:
                    print(entry)
                except:
                    pass
            else:
                outpath = 'E:/Golang/mw/catelogs/' + entry + '.html'
            allentry.add(entry)
            sum += 1
            if (sum % 900) == 0:
                newSession()
            with closing(rs.get(url, stream=True, timeout=30, proxies={"http": "http://{}".format(proxy)})) as r:
                # print(r.url)
                rc = r.status_code
                if 299 < rc or rc < 200:
                    print('returnCode%s\t%s' % (rc, url))
                    return
                ob = b''
                try:
                    with open(outpath, 'wb') as f:
                        for data in r.iter_content(1024):
                            f.write(data)
                            ob += data
                except:
                    with open('E:/Golang/mw/thsr/0error' + str(sum) + '.html', 'wb') as f:
                        for data in r.iter_content(1024):
                            f.write(data)
                            ob += data
            return ob
        except:
            print('what wrong?')
            delete_proxy(proxy)
            newSession()
            retry -= 1
    return False


if __name__ == '__main__':
    main()
