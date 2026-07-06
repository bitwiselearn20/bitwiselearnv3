// worker (email queue consumer) — Container App per plan section 6.2.
// Scales to zero: nothing in the queue means no cost.
param name string
param location string
param environmentId string
param acrLoginServer string
param acrUsername string
@secure()
param acrPassword string
param imageTag string

@secure()
param mqClient string
param emailUser string
@secure()
param emailPass string

resource worker 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  properties: {
    environmentId: environmentId
    configuration: {
      // No ingress: the worker only consumes from RabbitMQ, it never serves
      // HTTP traffic.
      registries: [
        {
          server: acrLoginServer
          username: acrUsername
          passwordSecretRef: 'acr-password'
        }
      ]
      secrets: [
        { name: 'acr-password', value: acrPassword }
        { name: 'mq-client', value: mqClient }
        { name: 'email-pass', value: emailPass }
      ]
    }
    template: {
      containers: [
        {
          name: 'worker'
          image: '${acrLoginServer}/bitwise-worker:${imageTag}'
          resources: {
            cpu: json('0.25')
            memory: '0.5Gi'
          }
          env: [
            { name: 'MQ_CLIENT', secretRef: 'mq-client' }
            { name: 'EMAIL_USER', value: emailUser }
            { name: 'EMAIL_PASS', secretRef: 'email-pass' }
          ]
        }
      ]
      scale: {
        minReplicas: 0
        maxReplicas: 10
        rules: [
          {
            name: 'queue-depth'
            custom: {
              type: 'rabbitmq'
              metadata: {
                queueName: 'email_jobs'
                mode: 'QueueLength'
                value: '20'
              }
              auth: [
                { secretRef: 'mq-client', triggerParameter: 'host' }
              ]
            }
          }
        ]
      }
    }
  }
}
