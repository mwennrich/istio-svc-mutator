apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: istio-svc-mutator
webhooks:
- admissionReviewVersions:
  - v1beta1
  clientConfig:
    caBundle: CA_BUNDLE
    service:
      name: istio-svc-mutator
      namespace: istio-ingress
      port: 443
  failurePolicy: Ignore
  matchPolicy: Exact
  name: istio-svc-mutator.metal-stack.dev
  namespaceSelector: {}
  objectSelector:
    matchExpressions:
    - key: app
      operator: In
      values:
      - istio-ingressgateway
  rules:
  - apiGroups:
    - ""
    apiVersions:
    - v1
    operations:
    - CREATE
    - UPDATE
    resources:
    - services
    scope: '*'
  sideEffects: None
