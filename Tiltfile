# Deploy: tell Tilt what YAML to deploy
k8s_yaml('build/deployment.yaml')

# Build: tell Tilt what images to build from which directories
docker_build('anthonycorbacho/chatbot', '.')

# Watch: tell Tilt how to connect locally (optional)
k8s_resource('chatbot-demo', port_forwards=[3000, 4000, 9000])