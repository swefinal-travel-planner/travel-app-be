pipeline {
    agent any
    environment {
        PORT =  credentials('BACKEND_PORT')

        DB_HOST = credentials('DB_HOST')
        DB_PORT = credentials('DB_PORT')
        DB_DATABASE = credentials('DB_NAME')
        DB_USERNAME = credentials('DB_USERNAME')
        DB_PASSWORD = credentials('DB_PASWORD') 
        DB_ROOT_PASSWORD = credentials('DB_ROOT_PASSWORD')

        REDIS_HOST = credentials('REDIS_HOST')
        REDIS_PORT = credentials('REDIS_PORT')

        JWT_SECRET = credentials('JWT_SECRET') 
    }

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
                    sh 'docker build \
                        -e PORT=$PORT \
                        -e DB_HOST=$DB_HOST \
                        -e DB_PORT=$DB_PORT \
                        -e DB_DATABASE=$DB_DATABASE \
                        -e DB_USERNAME=$DB_USERNAME \
                        -e DB_PASSWORD=$DB_PASSWORD \
                        -e DB_ROOT_PASSWORD=$DB_ROOT_PASSWORD \
                        -e REDIS_HOST=$REDIS_HOST \
                        -e REDIS_PORT=$REDIS_PORT \
                        -e JWT_SECRET=$JWT_SECRET \
                        -t travel-be .'
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
