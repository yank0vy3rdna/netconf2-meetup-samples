(sysrepoctl -l | grep loadbalancer) || (sysrepoctl -i yang/loadbalancer.yang -I yang/loadbalancer.xml)

netopeer2-server -d -v3 &
