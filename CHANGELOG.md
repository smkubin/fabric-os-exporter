## master / unreleased

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