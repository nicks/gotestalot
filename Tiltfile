# -*- mode: Python -*-

def main():
  yaml = local('cat deployments/gotestalot.yaml')
  image_name = 'cr.nick-santos.com/gotestalot'
  img = build_docker_image('Dockerfile.base', image_name, '/go/bin/gotestalot')
  img.add(local_git_repo('.'), '/go/src/github.com/nicks/gotestalot')
  img.add(local_git_repo('../../windmilleng/tilt/'), '/go/src/github.com/windmilleng/tilt')
  img.run('go install github.com/nicks/gotestalot')
  return k8s_service(yaml, img)
