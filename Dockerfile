FROM centos:7

EXPOSE 9090

ADD ./bin/shopify-customer-service /usr/bin/shopify-customer-service

CMD ["shopify-customer-service", "serve"]