pipeline {
    agent any

    stages {
        stage('Remove Old Docker Image') {
            steps {
                script {
                    // Remove the old image (optional)
                    // This assumes the old image is tagged with 'travel-be' and you want to remove it before re-building
                    sh 'docker rmi $(docker images -q travel-be) || true'  // || true to prevent error if no image exists
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Build the new Docker image
                    sh 'docker build -t travel-be .'
                }
            }
        }

        stage('Run Container') {
            steps {
                script {
                    // Run the newly built container
                    sh 'docker run -d --rm --name travel-be-container travel-be'
                }
            }
        }
    }
}
