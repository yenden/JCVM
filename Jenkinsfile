pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        bat 'go build ./core/method.go'
      }
    }
  }
}