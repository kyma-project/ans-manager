#!/bin/bash
NOTIFICATIONS_SB=${1:-stage-n-sb.json}
EVENTS_SB=${2:-stage-e-sb.json}
./token-n.sh ${NOTIFICATIONS_SB}
./token-e.sh ${EVENTS_SB} 
