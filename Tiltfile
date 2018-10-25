# -*- mode: Python -*-

def webpack():
  yaml = read_file('deployments/webpack.yaml')
  image_name = 'gcr.io/windmill-public-containers/gotestalot-webpack'
  entry = 'node_modules/.bin/webpack-dev-server --host 0.0.0.0 --port 8001'
  start_fast_build('Dockerfile.webpack.base', image_name, entry)
  repo = local_git_repo('.')
  add(repo.path('package.json'), '/app/package.json')
  add(repo.path('webpack.config.js'), '/app/webpack.config.js')
  add(repo.path('.babelrc'), '/app/.babelrc')
  add(repo.path('web/src'), '/app/web/src')

  run('npm install .', trigger='package.json')

  s = k8s_service(yaml, stop_build())
  s.port_forward(8001)
  return s

def server():
  yaml = read_file('deployments/server.yaml')
  image_name = 'gcr.io/windmill-public-containers/gotestalot'
  start_fast_build('Dockerfile.base', image_name, '/go/bin/gotestalot')
  repo = local_git_repo('.')
  add(repo, '/go/src/github.com/nicks/gotestalot')
  add(local_git_repo('../../windmilleng/tilt/'), '/go/src/github.com/windmilleng/tilt')
  run('go install github.com/nicks/gotestalot')

  s = k8s_service(yaml, stop_build())
  s.port_forward(8000)
  return s

def main():
  return composite_service([webpack, server])
