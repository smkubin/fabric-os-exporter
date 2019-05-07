FROM quay.io/prometheus/busybox:latest
# FROM s390x/busybox:latest
COPY fabric-os-exporter /bin/fabric-os-exporter
EXPOSE 9879
ENTRYPOINT ["/bin/fabric-os-exporter"]