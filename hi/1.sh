#!/bin/bash

API_SERVER="http://oswatcher-operator-controller-manager.oswatcher-operator-system.svc:8080/report"  # 여기에 컨트롤러가 띄워진 IP:Port 입력
NODE_NAME=$(hostname)
NODE_IP=$(hostname -I | awk '{print $1}')
COLLECTED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
UPTIME=$(uptime -p | cut -d' ' -f2-)
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'.' -f1)%
MEM_USAGE=$(free | awk '/Mem:/ {printf("%.0f%%", $3/$2 * 100)}')
ROOT_USAGE=$(df / | awk 'NR==2 {print $5}')

# Top CPU process
TOP_CPU_PROC=$(ps -eo pid,user,comm,%cpu,%mem --sort=-%cpu | awk 'NR==2 {printf("{\"pid\":\"%s\",\"user\":\"%s\",\"command\":\"%s\",\"cpu\":\"%s%%\",\"mem\":\"%s%%\"}", $1, $2, $3, $4, $5)}')

# Top MEM process
TOP_MEM_PROC=$(ps -eo pid,user,comm,%cpu,%mem --sort=-%mem | awk 'NR==2 {printf("{\"pid\":\"%s\",\"user\":\"%s\",\"command\":\"%s\",\"cpu\":\"%s%%\",\"mem\":\"%s%%\"}", $1, $2, $3, $4, $5)}')

JSON_PAYLOAD=$(cat <<EOF
{
  "apiVersion": "monitoring.oswatcher.io/v1",
  "kind": "OSStatus",
  "metadata": {
    "name": "$NODE_NAME",
    "namespace": "oswatcher-operator-system"
  },
  "spec": {
    "nodeName": "$NODE_NAME",
    "nodeIP": "$NODE_IP",
    "cpuUsage": "$CPU_USAGE",
    "memUsage": "$MEM_USAGE",
    "rootUsage": "$ROOT_USAGE",
    "uptime": "$UPTIME",
    "collectedAt": "$COLLECTED_AT",
    "topCPUProcs": [ $TOP_CPU_PROC ],
    "topMemProcs": [ $TOP_MEM_PROC ]
  }
}
EOF
)

curl -s -X POST -H "Content-Type: application/json" -d "$JSON_PAYLOAD" "$API_SERVER"
