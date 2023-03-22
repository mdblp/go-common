@Library('mdblp-library') _
def builderImage
pipeline {
    agent {
        label 'blp'
    }
    stages {
        stage('Initialization') {
            steps {
                script {
                    utils.initPipeline()
                    if(env.GIT_COMMIT == null) {
                        // git commit id must be a 40 characters length string (lower case or digits)
                        env.GIT_COMMIT = "f".multiply(40)
                    }
                    env.RUN_ID = UUID.randomUUID().toString()
                }
            }
        }
        stage('Build') {
            agent {
                docker {
                    image 'docker.ci.diabeloop.eu/go-build:1.17'
                    label 'blp'
                }
            }
            steps {
                script {
                    sh "cd ./v1 && go build && cd ../v2 && go build"
                }
            }
        }
        stage('Test V1') {
             steps {
                echo 'start mongo to serve as a testing db'
                sh 'docker network create gocommon${RUN_ID} && docker run --rm -d --net=gocommon${RUN_ID} --name=mongo4gocommon${RUN_ID} mongo:4.2'
                script {
                    docker.image('docker.ci.diabeloop.eu/go-build:1.17').inside("--net=gocommon${RUN_ID}") {
                        sh "TIDEPOOL_STORE_ADDRESSES=mongo4gocommon${RUN_ID}:27017  TIDEPOOL_STORE_DATABASE=gocommon_test cd $WORKSPACE/v1 && ./test.sh"
                    }
                }
            }
            post {
                always {
                    sh 'docker stop mongo4gocommon${RUN_ID} && docker network rm gocommon${RUN_ID}'
                    junit 'test-report.xml'
                    archiveArtifacts artifacts: 'coverage.html', allowEmptyArchive: true
                }
            }
        }
        stage('Test V2') {
             steps {
                echo 'start mongo to serve as a testing db'
                sh 'docker network create gocommon${RUN_ID} && docker run --rm -d --net=gocommon${RUN_ID} --name=mongo4gocommon${RUN_ID} mongo:4.2'
                script {
                    docker.image('docker.ci.diabeloop.eu/go-build:1.17').inside("--net=gocommon${RUN_ID}") {
                        sh "TIDEPOOL_STORE_ADDRESSES=mongo4gocommon${RUN_ID}:27017  TIDEPOOL_STORE_DATABASE=gocommon_test cd $WORKSPACE/v2 && ./test.sh"
                    }
                }
            }
            post {
                always {
                    sh 'docker stop mongo4gocommon${RUN_ID} && docker network rm gocommon${RUN_ID}'
                    junit 'test-report-v2.xml'
                    archiveArtifacts artifacts: 'coverage-v2.html', allowEmptyArchive: true
                }
            }
        }
        stage('Publish') {
            when { branch "master" }
            steps {
                publish()
            }
        }
    }
}
