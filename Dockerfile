FROM grafana/grafana:8.1.3

USER root
RUN wget https://github.com/org-insights/aws-s-3/releases/download/v1.0.0/itay-s3-datasource-1.0.0.zip \
 -O /var/lib/grafana/plugins/itay-s3-datasource-1.0.0.zip \
 && cd /var/lib/grafana/plugins \
 && unzip itay-s3-datasource-1.0.0.zip \
 && rm itay-s3-datasource-1.0.0.zip

# ADD dist /var/lib/grafana/plugins/aws-s-3/dist
# ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS orginsights-s3-datasource
# ENV GF_SERVER_HTTP_PORT 80

USER grafana