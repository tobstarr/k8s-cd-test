#!/usr/bin/env groovy

node {
  stage "SCM"
  env.IMAGE_REPO = "127.0.0.1:30015/tobstarr/k8s-cd-test"

  checkout scm

  stage "Build"
  sh 'bash ./k8s/trigger.sh ./k8s/build.tpl.yml'

  stage "Test"
  sh 'bash ./k8s/trigger.sh ./k8s/test.tpl.yml'

  stage "Test Integration"
  sh 'bash ./k8s/trigger.sh ./k8s/test.integration.tpl.yml'
}
