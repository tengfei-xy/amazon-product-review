# amazon-product-review

一种根据亚马逊的产品ID（asin）爬取评论的适用多平台的工具



## 使用方式

- mysql

  - 创建amazon_product_review用户
  
    ```sql
      # mysql 5.7.x
      GRANT ALL PRIVILEGES ON amazon_product_review.* to 'amazon_product_review'@'%'  identified by 'password';
      flush privileges;
      
      # mysql 8.x
      create user 'amazon_product_review'@'%' identified by 'password';
      GRANT ALL PRIVILEGES ON amazon.* TO 'amazon_product_review'@'%' WITH GRANT OPTION;
      flush privileges;
    ```
  
  
  - 根据[DLL.sql](https://raw.githubusercontent.com/tengfei-xy/amazon-product-review/main/sql/ddl.sql)文件进行初始化数据库
  
    ```
    # 执行后，将丢弃原表的数据
    bin/mysql -u root -p < ddl.sql
    ```

  - 填写配置文件的mysql部分
  - 将亚马逊的可用cookie写入数据库的cookie表的cookie字段中。
  
- 填写配置文件，剩余部分

- 启动程序

  ```bash
  ./amazon-product-review -c config.yaml -asin Bxxxxxxxxxx
  ```



## 使用说明

- 计算评论者的名称，评论时间，评论星级，标题，款式，评论过图片，评论主体内容得到校验值，每当评论的ID所对应的校验值发生变化时，表示评论有更新，将重新更新，以同步到mysql中。

- 未登录的cookie，无法获取完整的评论

- 仅在德国域名中测试完成， 其他地区未有测试