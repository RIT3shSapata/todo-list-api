FROM couchbase:enterprise-7.6.5
# RUN ["chmod", "+x", "/entrypoint.sh"]
COPY dockerfiles/couchbase/configure.sh /configure.sh
RUN ["chmod", "+x", "/configure.sh"]
CMD ["/configure.sh"]