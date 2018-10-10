#!/bin/bash

readonly CHART_NAME=stream-ingester
readonly CHART_DIR=./helm/$CHART_NAME

CONSUL_ADDR=${CONSUL_ADDR:=consul.internal.liveplanetstage.net:80}
WORKSPACE=${WORKSPACE:=staging}
PROJECT=${PROJECT:=cloud}
VERSION=${VERSION:=`git describe --abbrev=0`-`git rev-parse --short HEAD`}

function log {
  local readonly level="$1"
  local readonly message="$2"
  local readonly timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  >&2 echo -e "${timestamp} [${level}] [$SCRIPT_NAME] ${message}"
}

function log_info {
  local readonly message="$1"
  log "INFO" "$message"
}

function log_warn {
  local readonly message="$1"
  log "WARN" "$message"
}

function log_error {
  local readonly message="$1"
  log "ERROR" "$message"
}

function update_deps {
    log_info "Syncing dependencies..."
    helm dependencies update --kube-context ${KUBE_CONTEXT} ${CHART_DIR}
}

function has_jq {
  [ -n "$(command -v jq)" ]
}

function has_consul {
  [ -n "$(command -v consul)" ]
}

function has_helm {
  [ -n "$(command -v helm)" ]
}

function get_md {
  local readonly k="$1"
  local readonly v=`consul kv get -http-addr=${CONSUL_ADDR} $k`
  echo "$v"
}

function get_vars {
    log_info "Getting variables..."
    readonly GOOGLE_PROJECT=$(get_md infra/common/${WORKSPACE}/${PROJECT}/gcp_project)
    readonly KUBE_CONTEXT=$(get_md infra/common/${WORKSPACE}/${PROJECT}/kube_context)
    readonly INGRESS_STATIC_IP_NAME=$(get_md infra/services/stream-ingester/${WORKSPACE}/${PROJECT}/vars/ingress_static_ip_name)
    readonly INGRESS_HTTP_HOST=$(get_md infra/services/stream-ingester/${WORKSPACE}/${PROJECT}/vars/ingress_http_host)
    readonly RTMP_LB_IP=$(get_md infra/services/stream-ingester/${WORKSPACE}/${PROJECT}/vars/rtmp_lb_ip)
}

function deploy {
    log_info "Deploying ${CHART_NAME} version ${VERSION}"
    readonly REPOSITORY="us.gcr.io/${GOOGLE_PROJECT}/${CHART_NAME}"
    helm upgrade \
        --kube-context "${KUBE_CONTEXT}" \
        --install \
        --set image.tag="${VERSION}" \
        --set image.repository="${REPOSITORY}" \
        --set ingress.staticIpName="${INGRESS_STATIC_IP_NAME}" \
        --set ingress.httpHost="${INGRESS_HTTP_HOST}" \
        --set rtmpService.loadBalancerIP="${RTMP_LB_IP}" \
        --wait ${CHART_NAME} ${CHART_DIR}
}

if ! $(has_jq); then
    log_error "Could not find jq"
    exit 1
fi

if ! $(has_consul); then
    log_error "Could not find consul"
    exit 1
fi

if ! $(has_helm); then
    log_error "Could not find helm"
    exit 1
fi

get_vars
update_deps
deploy

exit $?