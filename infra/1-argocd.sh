#!/usr/bin/env bash
helm upgrade --install --create-namespace argocd argo/argo-cd -n argocd
