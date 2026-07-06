// RabbitMQ — single replica, ephemeral storage. Fine for the email/report
// job queues this system uses (jobs are re-triggerable — a lost email job
// just means "created" emails don't retroactively resend), but queued
// messages ARE lost on restart/redeploy since there's no persistent volume
// wired up here. For a production HA broker, point MQ_CLIENT at a managed
// service (CloudAMQP, Azure Service Bus premium tier) instead of this module.
param name string
param location string
param environmentId string

resource rabbitmq 'Microsoft.App/containerApps@2024-03-01' = {
  name: name
  location: location
  properties: {
    environmentId: environmentId
    configuration: {
      ingress: {
        external: false
        targetPort: 5672
        transport: 'tcp'
      }
    }
    template: {
      containers: [
        {
          name: 'rabbitmq'
          image: 'rabbitmq:3-management'
          resources: {
            cpu: json('0.5')
            memory: '1.0Gi'
          }
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 1
      }
    }
  }
}

output fqdn string = rabbitmq.properties.configuration.ingress.fqdn
