import requests
import time
import re
import dataset
import json
import logging
import argparse
import random
import zipfile
import os


from datetime import datetime
from selenium import webdriver
from selenium.webdriver import ChromeOptions
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.keys import Keys
from scrapy.selector import Selector
from pyvirtualdisplay import Display


class AmazonSearchResultParser:

    def __init__(self):
        # db = dataset.connect('sqlite:///amazon_search_result.db')
        db = dataset.connect('postgresql://postgres:postgres@localhost:5432/nichedb')
        # db = dataset.connect('mysql://user:password@localhost/mydatabase')
        self.table = db['amazon_search_result']

        SMARTPROXY_USER = 'spbb746f3f'
        SMARTPROXY_PASS = '9OTjKygF09Db'
        SMARTPROXY_URL = 'http://{}:{}@us.smartproxy.com:10000'.format(SMARTPROXY_USER, SMARTPROXY_PASS)
        self.proxies = {
          'http': SMARTPROXY_URL,
          'https': SMARTPROXY_URL,
        }

    def init_driver(self):
        print("1")
        self.create_auth_zip()
        print("2")

        chrome_options = ChromeOptions()
        chrome_options.add_extension('./DS Amazon Quick View.crx')
        print("3")

        # chrome_options.add_argument("user-data-dir=extentions")
        pluginfile = 'proxy_auth_plugin.zip'

        chrome_options.add_extension(pluginfile)
        print("4")

        path_to_chromedriver = './chromedriver'
        driver = webdriver.Chrome(options=chrome_options, executable_path=path_to_chromedriver)
        print("5")

        return driver

    def parse_keywords(self, keywords):
        try:
            display = Display(visible=0, size=(800, 600))
            display.start()
            driver = self.init_driver()

            for keyword in keywords:
                print(f'keyword = {keyword}')
                self.parse(keyword, driver)
        except Exception as e:
            print(f'ERROR - {repr(e)}')
        finally:
            driver.quit()
            display.stop()

    def parse(self, keyword, driver):
        try:
            driver.get(f'https://www.amazon.com/s?k={keyword}')
            time.sleep(2)
            print(f'get - https://www.amazon.com/s?k={keyword}')
            SCROLL_PAUSE_TIME = 3

            last_height = driver.execute_script("return document.body.scrollHeight")

            while True:
                for i in range(0, 5):
                    driver.find_element_by_tag_name('body').send_keys(Keys.PAGE_DOWN)
                time.sleep(SCROLL_PAUSE_TIME)
                new_height = driver.execute_script("return document.body.scrollHeight")

                if new_height == last_height:
                    break
                last_height = new_height

            time.sleep(2)

            response = Selector(text=driver.page_source)
            scraping_date = datetime.now()

            rows = []

            for item in response.xpath('//div[@data-asin and @data-component-type="s-search-result"]'):
                product_title = item.css('a.a-link-normal > span.a-text-normal::text').get()

                rating_count = item.xpath(".//a[contains(@href, 'Review')]/span/text()").get()
                rating_count = re.sub(r'\D', '', rating_count) if rating_count else None

                average_rating = item.css('i.a-icon-star-small > span::text').get()
                average_rating = average_rating.split('out')[0] if average_rating else average_rating
                price = item.css('span.a-price > span.a-offscreen::text').get()

                rank = item.css('span.extension-rank::text').get()
                rank = re.sub(r'\D', '', rank) if rank else rank

                category = item.xpath(".//span[contains(@class, 'extension-rank')]/parent::span/text()").get()
                category = re.sub(r'in|\(', '', category) if category else category

                asin = item.css('::attr(data-asin)').get()

                amazon_item = {
                    'scraping_date': str(scraping_date),
                    'query_id': keyword,
                    'asin': asin,
                    'product_title': product_title,
                    'rating_count': rating_count,
                    'average_rating': average_rating,
                    'price': price,
                    'rank': rank,
                    'category': category.strip() if category else category,
                    'sales': None
                }
                if rank and category:
                    amazon_item['sales'] = self.get_sales(rank, category)

                rows.append(amazon_item)
                self.table.insert(amazon_item)
                print(amazon_item)

        except Exception as e:
            print(f'ERROR - {repr(e)}')
        finally:
            pass

    def get_sales(self, rank, category):
        logging.info(f"checking sales for rank [{rank}] and category [{category}]")
        headers = {
            'authority': 'api.junglescout.com',
            'accept': 'application/json, text/javascript, */*; q=0.01',
            'user-agent': 'Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36',
            'origin': 'https://www.junglescout.com',
            'sec-fetch-site': 'same-site',
            'sec-fetch-mode': 'cors',
            'sec-fetch-dest': 'empty',
            'referer': 'https://www.junglescout.com/',
            'accept-language': 'en-US,en;q=0.9',
        }

        params = (
            ('rank', rank),
            ('category', category),
            ('store', 'us'),
        )
        try:
            response = requests.get(
                'https://api.junglescout.com/api/v1/sales_estimator',
                headers=headers,
                proxies=self.proxies,
                params=params)

            j_obj = response.json()

            if j_obj['status']:
                return j_obj['estSalesResult']
        except Exception as e:
            print(f'Cant get sales for {rank} - {category}, due to error {repr(e)}')

        return None

    def create_auth_zip(self):
        PROXY_HOST = 'us.smartproxy.com'
        PROXY_PORT = 10000
        PROXY_USER = 'spbb746f3f'
        PROXY_PASS = '9OTjKygF09Db'

        manifest_json = """
        {
            "version": "1.0.0",
            "manifest_version": 2,
            "name": "Chrome Proxy",
            "permissions": [
                "proxy",
                "tabs",
                "unlimitedStorage",
                "storage",
                "<all_urls>",
                "webRequest",
                "webRequestBlocking"
            ],
            "background": {
                "scripts": ["background.js"]
            },
            "minimum_chrome_version":"22.0.0"
        }
        """

        background_js = """
        var config = {
                mode: "fixed_servers",
                rules: {
                singleProxy: {
                    scheme: "http",
                    host: "%s",
                    port: parseInt(%s)
                },
                bypassList: ["localhost"]
                }
            };

        chrome.proxy.settings.set({value: config, scope: "regular"}, function() {});

        function callbackFn(details) {
            return {
                authCredentials: {
                    username: "%s",
                    password: "%s"
                }
            };
        }

        chrome.webRequest.onAuthRequired.addListener(
                    callbackFn,
                    {urls: ["<all_urls>"]},
                    ['blocking']
        );
        """ % (PROXY_HOST, PROXY_PORT, PROXY_USER, PROXY_PASS)

        pluginfile = 'proxy_auth_plugin.zip'
        if not os.path.exists(pluginfile):
            with zipfile.ZipFile(pluginfile, 'w') as zp:
                zp.writestr("manifest.json", manifest_json)
                zp.writestr("background.js", background_js)


if __name__ == '__main__':
    # ap = argparse.ArgumentParser()

    # ap.add_argument("-k", "--keyword", required=True, help='Keyword to search on amazon, text should be "here inside double quotes" ')
    # args = vars(ap.parse_args())

    parser = AmazonSearchResultParser()
    parser.parse_keywords(['ps4', 'ps4 game'])
    # rows = []
    # for row in parser.parse(args['keyword']):
    #     print(row)
