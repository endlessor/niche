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
        db = dataset.connect('postgresql://postgres:postgres@localhost:5432/nichedb')
        self.amazon_table = db['amazon_py']
        self.phrase_table = db['product_discoveries']
        self.mi_table = db['market_arguments']
        self.asin_table = db['asins']

        SMARTPROXY_USER = 'spbb746f3f'
        SMARTPROXY_PASS = '9OTjKygF09Db'
        SMARTPROXY_URL = 'http://{}:{}@us.smartproxy.com:10000'.format(SMARTPROXY_USER, SMARTPROXY_PASS)
        self.proxies = {
          'http': SMARTPROXY_URL,
          'https': SMARTPROXY_URL,
        }

    def init_driver(self):
        self.create_auth_zip()
        chrome_options = ChromeOptions()
        chrome_options.add_extension('./DS Amazon Quick View.crx')
        pluginfile = 'proxy_auth_plugin.zip'
        chrome_options.add_extension(pluginfile)
        path_to_chromedriver = './chromedriver'
        driver = webdriver.Chrome(options=chrome_options, executable_path=path_to_chromedriver)

        return driver

    def parse_keywords(self):
        try:
            display = Display(visible=0, size=(800, 600))
            display.start()
            driver = self.init_driver()

            pds = self.phrase_table.all()

            for pd in pds:
                keyword = pd['phrase']
                preset = pd['preset']
                print(f'keyword = {keyword}, preset = {preset}')
                self.parse(driver, keyword, preset)
        except Exception as e:
            print(f'ERROR - {repr(e)}')
        finally:
            driver.quit()
            display.stop()

    def parse(self, driver, keyword, preset):
        try:
            scraping_date = datetime.now()
            db_date = str(scraping_date)[0:7]
            ex_phrase = self.asin_table.find_one(date=db_date, phrase=keyword)
            if ex_phrase:
                print("existing keyword:", keyword)
                return

            driver.get(f'https://www.amazon.com/s?k={keyword}')
            time.sleep(2)
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
                    'phrase': keyword,
                    'preset': preset,
                    'asin': asin,
                    'product_title': product_title,
                    'rating_count': rating_count,
                    'average_rating': average_rating,
                    'price': price,
                    'rank': rank,
                    'category': category.strip() if category else category,
                    'sales': None
                }
                # if rank and category and len(category) > 0:
                #     sales = self.get_sales(rank, category)
                #     if not str(sales).isdigit():
                #         sales = None
                #     amazon_item['sales'] = sales

                rows.append(amazon_item)
                self.amazon_table.insert(amazon_item)
                print("asin:", asin)
                # self.get_mi_asin(preset, keyword, asin)

            asin_data = {
                'date': db_date,
                'phrase': keyword
            }
            self.asin_table.insert(asin_data)

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

    def get_mi_asin(self, preset, keyword, asin):
        searchListingURL = "https://viral-launch.com/sellers/assets/php/market-intelligence/search-listing.php"
        headers = {
            "Content-Type": "application/json",
        }
        params = (
            ("asin", asin),
        )
        body = {
            'by': "subs@coral8.co",
            'marketplace': "US",
            'objectId': "kN1G4gsZEh",
            'phrase': keyword,
            'source': "viral-launch.com"
        }
        try:
            response = requests.post(
                searchListingURL,
                json=body,
                headers=headers,
                params=params)

            result = response.json()

            if not result.get('title'):
                return
            mi_item = {
                'phrase': keyword,
                'preset': preset,
                'asin': asin,
                'at': result.get('at', ""),
                'title': result.get('title', ""),
                'description': result.get('description', ""),
                'price': result.get('price', 0),
                'parent_asin': result.get('parentASIN', ""),
                'brand': result.get('brand', ""),
                'bsr': result.get('bsr', 0),
                'category': result.get('category', ""),
                'category_id': result.get('categoryID', 0),
                'categories': result.get('categories', []),
                'features': result.get('features', []),                     
                'image_urls': result.get('imageUrls', []),                    
                'net_profit': result.get('netProfit', 0),                    
                'revenue': result.get('revenue', 0),                      
                'review_count': result.get('reviewCount', 0),        
                'review_rate': result.get('reviewRate', 0),                   
                'review_rating': result.get('reviewRating', 0),      
                'price_amazon': result.get('priceAmazon', 0),
                'price_new': result.get('priceNew', 0),
                'product_group': result.get('productGroup', ""),
                'profit_margin': result.get('profitMargin', 0),                 
                'sales': result.get('sales', 0),              
                'sales_last_year': result.get('salesLastYear', 0),      
                'sales_to_reviews': result.get('salesToReviews', 0),               
                'seller_count': result.get('sellerCount', 0),        
                'unit_margin': result.get('unitMargin', 0),                   
                'star_rating': result.get('starRating', 0),                   
                'best_sales_period': result.get('bestSalesPeriod', ""),              
                'is_name_brand': result.get('isNameBrand', False),                  
                'price_change_last_ninety_days': result.get('priceChangeLastNinetyDays', 0),    
                'review_count_change_monthly': result.get('reviewCountChangeMonthly', 0),     
                'sales_change_last_ninety_days': result.get('salesChangeLastNinetyDays', 0),    
                'sales_pattern': result.get('salesPattern', ""),
                'sales_year_over_year': result.get('salesYearOverYear', 0),            
                'initial_cost': result.get('initialCost', 0),                  
                'initial_net_profit': result.get('initialNetProfit', 0),             
                'initial_organic_sales_projection': result.get('initialOrganicSalesProjection', 0),
                'initial_units_to_order': result.get('initialUnitsToOrder', 0),
                'ongoing_organic_sales_projection': result.get('ongoingOrganicSalesProjection', 0),
                'ongoing_units_to_order': result.get('ongoingUnitsToOrder', 0),
                'promotion_duration': result.get('promotionDuration', 0),  
                'promotion_units_daily': result.get('promotionUnitsDaily', 0),
                'fulfillment': result.get('fulfillment', ""), 
                'promotion_units_total': result.get('promotionUnitsTotal', 0),
                'is_variation_with_shared_bsr': result.get('isVariationWithSharedBSR', False),     
                'offer_count_new': result.get('offerCountNew', 0),      
                'offer_count_used': result.get('offerCountUsed', 0),     
                'package_height': result.get('packageHeight', 0),                
                'package_length': result.get('packageLength', 0),                
                'package_quantity': result.get('packageQuantity', 0),
                'package_weight': result.get('packageWeight', 0),
                'package_width': result.get('packageWidth', 0),
            }
            self.mi_table.insert(mi_item)
        except Exception as e:
            print(f'Cant get mi, due to error {repr(e)}')
        return 

if __name__ == '__main__':
    parser = AmazonSearchResultParser()
    parser.parse_keywords()
