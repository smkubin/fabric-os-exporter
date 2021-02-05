## 0.5.3 / 2021-02-05
### Changes
* [CHANGE] Have the metric name compatible with fec_cor_detected and fec_corrected_rate
* [FIXBUG] Fix parsing error when the fec_cor_detected metric is missing

## 0.5.2 / 2020-11-23
### Changes
* [CHANGE] Fix CrossSiteRequestForgery vulnerabilities

## 0.5.1 / 2020-01-17
### Changes
* [CHANGE] Use go module to organize the dependency modules

## 0.5.0 / 2019-09-11
### Changes
* [CHANGE] Diable http methods other than GET
* [FIXBUG] Handle err when the sanswitch is disabled

## 0.4.0 / 2019-07-10
### Changes
* [FIX] Fix the 'uptime' metric
* [FEATURE] Add more debug information
* [FIX] Fix the error that temperature and fan_speed have wrong label
* [CHANGE] Log more information in error path

## 0.3.0 / 2019-06-28
### Changes
* [FEATURE] Add 'target' label, its value is ipaddress
* [CHANGE] Add enable-full-metrics flag to enable advanced metrics

## 0.2.0 / 2019-06-15
### **Breaking changes**
* [CHANGE] Change method of reading username and password from cli to config file.
* [FEATURE] Add metric of FEC corrected for fabric-os-exporter
### Changes
* [FEATURE] Add label of version for uptime metric
* [CHANGE] Change the label 'target' from ip into hostname
* [CHANGE] Change the label name from 'target' into 'resource'
* [CHANGE] Fix import statement
* [CHANGE] Fix Dockerfile

## 0.1.0 / 2019-04-12
* [CLEANUP] Introduced semantic versioning and changelog. From now on,
  changes will be reported in this file.
* [ENHANCEMENT] First version collects metrics from CLI command uptime,
  porterrshow, sensorshow.
