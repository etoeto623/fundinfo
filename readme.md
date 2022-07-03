> 本项目旨在爬取基金的一些信息，用来进行投资参考
# 基金净值
- 新浪财经数据源  
基金净值使用的接口是新浪财经的：[https://stock.finance.sina.com.cn/fundInfo/api/openapi.php/CaihuiFundInfoService.getNav?symbol=000961&datefrom=&dateto=&page=4](https://stock.finance.sina.com.cn/fundInfo/api/openapi.php/CaihuiFundInfoService.getNav?symbol=000961&datefrom=&dateto=&page=4)
- 天天基金数据源   
天天基金的页面是：[http://fundf10.eastmoney.com/jjjz_013146.html](http://fundf10.eastmoney.com/jjjz_013146.html)，对应查询数据查询接口是：[http://api.fund.eastmoney.com/f10/lsjz?callback=jQuery183003819373546541671_1657337229493&fundCode=013146&pageIndex=1&pageSize=20&startDate=&endDate=&_=1657337229609](http://api.fund.eastmoney.com/f10/lsjz?callback=jQuery183003819373546541671_1657337229493&fundCode=013146&pageIndex=1&pageSize=20&startDate=&endDate=&_=1657337229609)
# roadmap
- [X] 拉取基金净值数据
- [ ] 拉取基金基本信息