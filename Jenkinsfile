pipeline {
    agent any
    environment {
        PORT =  credentials('BACKEND_PORT')

        DB_HOST = credentials('DB_HOST')
        DB_PORT = credentials('DB_PORT')
        DB_DATABASE = credentials('DB_DATABASE')
        DB_USERNAME = credentials('DB_USERNAME')
        DB_PASSWORD = credentials('DB_PASSWORD') 
        DB_ROOT_PASSWORD = credentials('DB_ROOT_PASSWORD')

        REDIS_HOST = credentials('REDIS_HOST')
        REDIS_PORT = credentials('REDIS_PORT')

        JWT_SECRET = credentials('JWT_SECRET') 
    }

    stages {
        stage('Remove Old Docker Image') {
            steps {
                script {
                    // Stop and remove the old container if it exists
                    sh 'docker stop travel-be-container |   | true'
                    sh 'docker rm travel-be-container || true'
                    
                    // Remove the old image
                    sh 'docker rmi $(docker images -q travel-be) || true'  // || true to prevent error if no image exists
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Build the new Docker image
                    sh 'docker build \
                        --build-arg PORT=$PORT \
                        --build-arg DB_HOST=$DB_HOST \
                        --build-arg DB_PORT=$DB_PORT \
                        --build-arg DB_DATABASE=$DB_DATABASE \
                        --build-arg DB_USERNAME=$DB_USERNAME \
                        --build-arg DB_PASSWORD=$DB_PASSWORD \
                        --build-arg DB_ROOT_PASSWORD=$DB_ROOT_PASSWORD \
                        --build-arg REDIS_HOST=$REDIS_HOST \
                        --build-arg REDIS_PORT=$REDIS_PORT \
                        --build-arg JWT_SECRET=$JWT_SECRET \
                        -t travel-be .'
                }
            }
        }

        stage('Run Container') {
            steps {
                script {
                    // Run the newly built container with restart policy
                    sh 'docker run --rm -d --restart unless-stopped --name travel-be-container -p 9090:9090 travel-be'
                }
            }
        }
    }
}
