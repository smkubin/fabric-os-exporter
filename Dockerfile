FROM s390x/busybox:latest
# FROM quay.io/prometheus/busybox:latest
COPY fabricos.yaml /etc/fabric-os-exporter/fabricos.yaml
COPY fabric-os-exporter /bin/fabric-os-exporter
EXPOSE 9879
ENTRYPOINT ["/bin/fabric-os-exporter"]
CMD ["--config.file=/etc/fabric-os-exporter/fabricos.yaml"]
