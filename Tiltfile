# -*- mode: Python -*-

def main():
  yaml = local('cat deployments/gotestalot.yaml')
  image_name = 'gcr.io/windmill-public-containers/gotestalot'
  entry = '/go/bin/gotestalot --web_dir=/go/src/github.com/nicks/gotestalot/web github.com/windmilleng/tilt'
  img = build_docker_image('Dockerfile.base', image_name, entry)
  img.add(local_git_repo('.'), '/go/src/github.com/nicks/gotestalot')
  img.add(local_git_repo('../../windmilleng/tilt/'), '/go/src/github.com/windmilleng/tilt')
  img.run('go install github.com/nicks/gotestalot')
  return k8s_service(yaml, img)
